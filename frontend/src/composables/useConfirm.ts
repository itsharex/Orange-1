/**
 * @file composables/useConfirm.ts
 * @description 全局确认框 Hook
 * 提供函数调用的方式唤起确认模态框，基于 Vue createApp 动态挂载组件。
 */
import { ref, createApp, h } from 'vue'
import ConfirmModal from '@/components/common/ConfirmModal.vue'

interface ConfirmOptions {
  title?: string   // 标题
  message: string  // 内容消息
}

// 模态框实例接口
interface ModalInstance {
  modalRef: {
    open: (options: { title?: string; message: string }) => Promise<boolean>
  }
}

// 单例模式：全局共享一个模态框实例
let modalInstance: ModalInstance | null = null
let mountNode: HTMLElement | null = null

/** 创建或获取全局模态框实例 */
const createModalInstance = (): ModalInstance => {
  if (modalInstance) return modalInstance

  // 创建 DOM 容器
  mountNode = document.createElement('div')
  document.body.appendChild(mountNode)

  // 动态创建 Vue 应用实例
  const app = createApp({
    setup() {
      const modalRef = ref<InstanceType<typeof ConfirmModal> | null>(null)
      return { modalRef }
    },
    render() {
      return h(ConfirmModal, { ref: 'modalRef' })
    },
  })

  // 挂载
  modalInstance = app.mount(mountNode) as unknown as ModalInstance
  return modalInstance
}

/**
 * 使用确认框
 * @returns { confirm } 函数
 */
export const useConfirm = () => {
  /**
   * 弹出确认框
   * @param options 消息字符串或配置对象
   * @returns Promise<boolean> 确认返回 true, 取消返回 false
   */
  const confirm = async (options: ConfirmOptions | string): Promise<boolean> => {
    const instance = createModalInstance()

    // 等待 Next Tick 确保组件挂载完成且 ref 可用
    await new Promise((resolve) => setTimeout(resolve, 0))

    const opts = typeof options === 'string' ? { message: options } : options

    return instance.modalRef.open(opts)
  }

  return { confirm }
}

export default useConfirm
