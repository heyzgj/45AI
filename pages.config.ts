import { defineUniPages } from '@uni-helper/vite-plugin-uni-pages'

export default defineUniPages({
  globalStyle: {
    navigationStyle: 'default',
    navigationBarTitleText: '45AI',
    navigationBarBackgroundColor: '#FCFBF9',
    navigationBarTextStyle: 'black',
    backgroundColor: '#FCFBF9',
  },
  easycom: {
    autoscan: true,
    custom: {
      '^fg-(.*)': '@/components/fg-$1/fg-$1.vue',
      '^wd-(.*)': 'wot-design-uni/components/wd-$1/wd-$1.vue',
      '^(?!z-paging-refresh|z-paging-load-more)z-paging(.*)':
        'z-paging/components/z-paging$1/z-paging$1.vue',
    },
  },
  // 如果不需要tabBar，推荐使用 spa 模板。（pnpm create xxx -t spa）
  tabBar: {
    color: '#9B9B9B',
    selectedColor: '#E89B93',
    backgroundColor: '#FCFBF9',
    borderStyle: 'white',
    height: '56px',
    fontSize: '14px',
    list: [
      {
        pagePath: 'pages/gallery/index',
        text: '首页',
      },
      {
        pagePath: 'pages/profile/index',
        text: '我的',
      },
    ],
  },
})
