App({
  globalData: {
    // 后端 API 根地址（替换为真实域名）
    API_BASE: 'https://koopli.shop',
  },

  onLaunch() {
    // 初始化历史记录（如不存在则创建空数组）
    if (!wx.getStorageSync('history')) {
      wx.setStorageSync('history', []);
    }
  },
});
