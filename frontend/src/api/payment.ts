import { post } from '@/utils/request'

// Credit package definitions
export const getCreditPackages = () => {
  return [
    {
      id: 'pack_small',
      title: '50胶卷',
      description: '适合轻度使用',
      credits: 50,
      amount: 9.9,
      popular: false,
    },
    {
      id: 'pack_medium',
      title: '120胶卷',
      description: '性价比首选',
      credits: 120,
      amount: 19.9,
      originalPrice: 24.9,
      popular: true,
    },
    {
      id: 'pack_large',
      title: '300胶卷',
      description: '重度用户推荐',
      credits: 300,
      amount: 39.9,
      originalPrice: 49.9,
      popular: false,
    },
    {
      id: 'pack_extra',
      title: '600胶卷',
      description: '超值大包装',
      credits: 600,
      amount: 69.9,
      originalPrice: 99.9,
      popular: false,
    },
  ]
}

// Apple IAP product definitions
export const getAppleIAPProducts = () => {
  return [
    {
      id: 'com.45ai.credits.50',
      product_id: 'com.45ai.credits.50',
      title: '50胶卷',
      description: '适合轻度使用',
      credits: 50,
      price: '$0.99',
      popular: false,
    },
    {
      id: 'com.45ai.credits.120',
      product_id: 'com.45ai.credits.120',
      title: '120胶卷',
      description: '性价比首选',
      credits: 120,
      price: '$2.99',
      originalPrice: '$3.99',
      popular: true,
    },
    {
      id: 'com.45ai.credits.300',
      product_id: 'com.45ai.credits.300',
      title: '300胶卷',
      description: '重度用户推荐',
      credits: 300,
      price: '$6.99',
      originalPrice: '$9.99',
      popular: false,
    },
    {
      id: 'com.45ai.credits.600',
      product_id: 'com.45ai.credits.600',
      title: '600胶卷',
      description: '超值大包装',
      credits: 600,
      price: '$12.99',
      originalPrice: '$19.99',
      popular: false,
    },
  ]
}

// WeChat Pay integration
export const completePayment = async (params: { amount: number; description: string }) => {
  try {
    const response = await post('/billing/purchase', {
      amount: params.amount,
      description: params.description,
      platform: 'wechat',
    })
    return response
  } catch (error) {
    console.error('WeChat payment failed:', error)
    throw error
  }
}

// Apple IAP integration
export const completeApplePurchase = async (params: {
  product_id: string
  title: string
  description: string
}) => {
  try {
    const response = await post('/billing/purchase', {
      product_id: params.product_id,
      title: params.title,
      description: params.description,
      platform: 'apple',
    })
    return response
  } catch (error) {
    console.error('Apple IAP failed:', error)
    throw error
  }
}

// Price formatting utilities
export const formatPrice = (price: number): string => {
  return `¥${price.toFixed(2)}`
}

export const formatUSDPrice = (price: string): string => {
  return price // Already formatted as $X.XX
}
