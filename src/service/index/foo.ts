import { get, post } from '@/utils/request'
export interface IFooItem {
  id: string
  name: string
}

/** GET 请求 */
export const getFooAPI = (name: string) => {
  return get<IFooItem>('/foo', { name })
}
/** GET 请求；支持 传递 header 的范例 */
export const getFooAPI2 = (name: string) => {
  return get<IFooItem>(
    '/foo',
    { name },
    {
      header: { 'Content-Type-100': '100' },
    },
  )
}

/** POST 请求 */
export const postFooAPI = (name: string) => {
  return post<IFooItem>('/foo', { name })
}
/** POST 请求；需要传递 query 参数的范例；微信小程序经常有同时需要query参数和body参数的场景 */
export const postFooAPI2 = (name: string) => {
  return post<IFooItem>('/foo', { name })
}
/** POST 请求；支持 传递 header 的范例 */
export const postFooAPI3 = (name: string) => {
  return post<IFooItem>(
    '/foo',
    { name },
    {
      header: { 'Content-Type-100': '100' },
    },
  )
}
