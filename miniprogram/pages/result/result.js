Page({
  data: {
    imagePath: '',
    originPath: '',
    statusMsg: '',
  },

  onLoad(options) {
    const imagePath = decodeURIComponent(options.imagePath || '');
    const originPath = decodeURIComponent(options.originPath || '');
    this.setData({ imagePath, originPath });
    wx.setNavigationBarTitle({ title: 'ASCII 结果' });
  },

  saveToAlbum() {
    if (!this.data.imagePath) return;

    wx.saveImageToPhotosAlbum({
      filePath: this.data.imagePath,
      success: () => {
        this.setData({ statusMsg: '✅ 已保存到相册' });
        setTimeout(() => this.setData({ statusMsg: '' }), 2000);
      },
      fail: (err) => {
        if (err.errMsg && err.errMsg.includes('auth deny')) {
          wx.showModal({
            title: '需要权限',
            content: '请在设置中开启相册写入权限',
            confirmText: '去设置',
            success: (res) => {
              if (res.confirm) wx.openSetting();
            },
          });
        } else {
          this.setData({ statusMsg: '❌ 保存失败' });
        }
      },
    });
  },

  shareToFriend() {
    // 触发用户分享（onShareAppMessage 处理）
    wx.showShareMenu({
      withShareTicket: true,
      menus: ['shareAppMessage', 'shareTimeline'],
    });
    this.setData({ statusMsg: '长按图片可直接分享' });
  },

  onShareAppMessage() {
    return {
      title: '我用 ASCII 字符画了一张图！',
      imageUrl: this.data.imagePath,
    };
  },

  onShareTimeline() {
    return {
      title: 'ASCII 字符画转换',
      imageUrl: this.data.imagePath,
    };
  },
});
