const app = getApp();

Page({
  data: {
    imagePath: '',         // 本地临时文件路径
    previewSrc: '',        // 预览图路径
    converting: false,
    suggesting: false,
    suggestAvailable: false,
    statusMsg: '',
    params: {
      width: 100,
      colorful: false,
      chars: '@#S%?*+;:,. ',
    },
    charsetOptions: [
      { label: '默认',          value: '@#S%?*+;:,. ' },
      { label: '简洁',          value: '@%#*+=-:. ' },
      { label: '精细（Unicode）', value: '█▓▒░▐▌▄▀+;:,. ' },
      { label: '白底黑字',       value: ' .:-=+*#%@' },
    ],
    charsetIndex: 0,
  },

  onLoad() {
    // 检测 suggest 是否可用（服务端无 API Key 时返回 503）
    wx.request({
      url: app.globalData.API_BASE + '/convert/suggest',
      method: 'POST',
      success: (res) => {
        if (res.statusCode !== 503) {
          this.setData({ suggestAvailable: true });
        }
      },
      fail: () => {},
    });
  },

  chooseImage() {
    wx.chooseMedia({
      count: 1,
      mediaType: ['image'],
      sourceType: ['album', 'camera'],
      success: (res) => {
        const path = res.tempFiles[0].tempFilePath;
        this.setData({ imagePath: path, previewSrc: path, statusMsg: '' });
      },
    });
  },

  onWidthChange(e) {
    this.setData({ 'params.width': e.detail.value });
  },

  onCharsetChange(e) {
    const idx = e.detail.value;
    this.setData({
      charsetIndex: idx,
      'params.chars': this.data.charsetOptions[idx].value,
    });
  },

  onColorfulChange(e) {
    this.setData({ 'params.colorful': e.detail.value });
  },

  doSuggest() {
    if (!this.data.imagePath) return;
    this.setData({ suggesting: true, statusMsg: '' });

    wx.uploadFile({
      url: app.globalData.API_BASE + '/convert/suggest',
      filePath: this.data.imagePath,
      name: 'image',
      success: (res) => {
        if (res.statusCode !== 200) {
          this.setData({ statusMsg: '❌ 推荐失败' });
          return;
        }
        const params = JSON.parse(res.data);
        // 匹配字符集选项
        let charsetIndex = 0;
        const opts = this.data.charsetOptions;
        for (let i = 0; i < opts.length; i++) {
          if (opts[i].value === params.chars) {
            charsetIndex = i;
            break;
          }
        }
        this.setData({
          'params.width': params.width || 100,
          'params.colorful': params.colorful || false,
          'params.chars': params.chars || '@#S%?*+;:,. ',
          charsetIndex,
          statusMsg: `✨ 推荐：${params.reason || '参数已更新'}`,
        });
      },
      fail: () => {
        this.setData({ statusMsg: '❌ 网络错误' });
      },
      complete: () => {
        this.setData({ suggesting: false });
      },
    });
  },

  doConvert() {
    if (!this.data.imagePath || this.data.converting) return;
    this.setData({ converting: true, statusMsg: '转换中…' });

    const { width, colorful, chars } = this.data.params;

    wx.uploadFile({
      url: app.globalData.API_BASE + '/convert',
      filePath: this.data.imagePath,
      name: 'image',
      formData: {
        width: String(width),
        colorful: colorful ? 'true' : 'false',
        chars,
      },
      success: (res) => {
        if (res.statusCode !== 200) {
          this.setData({ statusMsg: '❌ 转换失败，请重试' });
          return;
        }
        // 将后端返回的 PNG 保存为临时文件
        const tmpPath = `${wx.env.USER_DATA_PATH}/ascii_${Date.now()}.png`;
        const fs = wx.getFileSystemManager();
        fs.writeFile({
          filePath: tmpPath,
          data: res.data,
          encoding: 'binary',
          success: () => {
            // 保存历史记录
            this._saveHistory(tmpPath);
            // 跳转到结果页
            wx.navigateTo({
              url: `/pages/result/result?imagePath=${encodeURIComponent(tmpPath)}&originPath=${encodeURIComponent(this.data.imagePath)}`,
            });
            this.setData({ statusMsg: '' });
          },
          fail: () => {
            this.setData({ statusMsg: '❌ 保存文件失败' });
          },
        });
      },
      fail: () => {
        this.setData({ statusMsg: '❌ 网络连接失败' });
      },
      complete: () => {
        this.setData({ converting: false });
      },
    });
  },

  _saveHistory(resultPath) {
    const history = wx.getStorageSync('history') || [];
    history.unshift({
      id: Date.now(),
      resultPath,
      originPath: this.data.imagePath,
      params: { ...this.data.params },
      time: new Date().toLocaleString('zh-CN'),
    });
    // 最多保留20条
    if (history.length > 20) history.splice(20);
    wx.setStorageSync('history', history);
  },
});
