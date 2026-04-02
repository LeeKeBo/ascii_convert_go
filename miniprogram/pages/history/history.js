Page({
  data: {
    history: [],
  },

  onShow() {
    // 每次显示时刷新列表
    const history = wx.getStorageSync('history') || [];
    this.setData({ history });
  },

  goToResult(e) {
    const item = e.currentTarget.dataset.item;
    wx.navigateTo({
      url: `/pages/result/result?imagePath=${encodeURIComponent(item.resultPath)}&originPath=${encodeURIComponent(item.originPath || '')}`,
    });
  },

  onImgError(e) {
    // 图片加载失败时（临时文件已过期）从历史中移除该条
    const idx = e.currentTarget.dataset.index;
    const history = wx.getStorageSync('history') || [];
    history.splice(idx, 1);
    wx.setStorageSync('history', history);
    this.setData({ history });
  },

  clearHistory() {
    wx.showModal({
      title: '确认清空',
      content: '将删除所有历史记录，不可恢复',
      success: (res) => {
        if (res.confirm) {
          wx.setStorageSync('history', []);
          this.setData({ history: [] });
        }
      },
    });
  },
});
