import { ref, createApp, h } from 'vue'
import ConfirmModal from '@/components/common/ConfirmModal.vue'

interface ConfirmOptions {
  title?: string
  message: string
}

interface ModalInstance {
  modalRef: {
    open: (options: { title?: string; message: string }) => Promise<boolean>
  }
}

let modalInstance: ModalInstance | null = null
let mountNode: HTMLElement | null = null

const createModalInstance = (): ModalInstance => {
  if (modalInstance) return modalInstance

  mountNode = document.createElement('div')
  document.body.appendChild(mountNode)

  const app = createApp({
    setup() {
      const modalRef = ref<InstanceType<typeof ConfirmModal> | null>(null)
      return { modalRef }
    },
    render() {
      return h(ConfirmModal, { ref: 'modalRef' })
    },
  })

  modalInstance = app.mount(mountNode) as unknown as ModalInstance
  return modalInstance
}

export const useConfirm = () => {
  const confirm = async (options: ConfirmOptions | string): Promise<boolean> => {
    const instance = createModalInstance()

    // Wait for next tick to ensure ref is available
    await new Promise((resolve) => setTimeout(resolve, 0))

    const opts = typeof options === 'string' ? { message: options } : options

    return instance.modalRef.open(opts)
  }

  return { confirm }
}

export default useConfirm
