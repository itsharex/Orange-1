<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { dictionaryApi, type Dictionary, type DictionaryItem } from '@/api/dictionary'
import { useToast } from '@/composables/useToast'
import { useConfirm } from '@/composables/useConfirm'

const toast = useToast()
const { confirm } = useConfirm()

// ============ 字典管理逻辑 ============
const activeDictId = ref<string>('') // 选中的字典 Code
const dictionaries = ref<Dictionary[]>([])
const activeDictItems = ref<DictionaryItem[]>([])

// Fetch dictionary list
const fetchDictionaries = async () => {
  try {
    const res = await dictionaryApi.list()
    if (res.data.code === 0) {
      dictionaries.value = res.data.data
      if (dictionaries.value.length > 0 && !activeDictId.value) {
        const firstDict = dictionaries.value[0]
        if (firstDict) {
          activeDictId.value = firstDict.code
        }
      }
    }
  } catch (error) {
    console.error('Failed to fetch dictionaries:', error)
  }
}

// Fetch items for selected dictionary
const fetchDictItems = async (code: string) => {
  if (!code) return
  try {
    const res = await dictionaryApi.getItems(code)
    if (res.data.code === 0) {
      activeDictItems.value = res.data.data
    }
  } catch (error) {
    console.error(`Failed to fetch items for ${code}:`, error)
  }
}

// Watch active dictionary selection change
watch(activeDictId, (newCode) => {
  if (newCode) {
    fetchDictItems(newCode)
  }
})

// 获取当前选中字典的名称
const activeDictName = computed(() => {
  const dict = dictionaries.value.find(d => d.code === activeDictId.value)
  return dict?.name || activeDictId.value
})

// Modal Logic
const showModal = ref(false)
const isEditing = ref(false)
const modalForm = ref({
  id: 0,
  label: '',
  value: '',
  sort: 0
})

const openAddModal = () => {
  isEditing.value = false
  modalForm.value = { id: 0, label: '', value: '', sort: activeDictItems.value.length + 1 }
  showModal.value = true
}

const openEditModal = (item: DictionaryItem) => {
  isEditing.value = true
  modalForm.value = {
    id: item.id,
    label: item.label,
    value: item.value,
    sort: item.sort
  }
  showModal.value = true
}

const handleModalSubmit = async () => {
  if (!activeDictId.value) return
  
  const label = modalForm.value.label.trim()
  const value = modalForm.value.value.trim()
  
  if (!label || !value) {
    toast.warning('请输入名称和值')
    return
  }
  
  try {
    let res
    if (isEditing.value) {
      res = await dictionaryApi.updateItem(activeDictId.value, modalForm.value.id, {
        label,
        value,
        sort: modalForm.value.sort
      })
    } else {
      res = await dictionaryApi.createItem(activeDictId.value, {
        label,
        value,
        sort: modalForm.value.sort
      })
    }
    
    if (res.data.code === 0) {
      toast.success(isEditing.value ? '修改成功' : '添加成功')
      showModal.value = false
      fetchDictItems(activeDictId.value)
    } else {
      toast.error(`${isEditing.value ? '修改' : '添加'}失败: ${res.data.message}`)
    }
  } catch (error) {
    console.error('Failed to save item:', error)
    toast.error(isEditing.value ? '修改失败' : '添加失败')
  }
}

const deleteDictItem = async (id: number) => {
  if (!activeDictId.value) return
  const confirmed = await confirm('确定要删除这个选项吗？')
  if (confirmed) {
    try {
      const res = await dictionaryApi.deleteItem(activeDictId.value, id)
      if (res.data.code === 0) {
        await fetchDictItems(activeDictId.value)
        toast.success('删除成功')
      }
    } catch (error) {
      console.error('Failed to delete item:', error)
      toast.error('删除失败')
    }
  }
}

onMounted(() => {
  fetchDictionaries()
})
</script>

<template>
  <div class="dict-management">
    <!-- 头部区域 -->
    <div class="dev-header">
      <div class="dev-header-content">
        <div class="dev-title-section">
          <div class="dev-icon-wrapper">
            <i class="ri-book-2-line"></i>
          </div>
          <div class="dev-title-info">
            <h2 class="dev-title">字典管理</h2>
            <p class="dev-subtitle">管理系统数据字典和配置项</p>
          </div>
        </div>
        <button class="dev-create-btn" @click="openAddModal">
          <i class="ri-add-line"></i>
          <span>新增条目</span>
        </button>
      </div>

      <!-- 统计卡片 -->
      <div class="dev-stats">
        <div class="dev-stat-card">
          <div class="dev-stat-icon total">
            <i class="ri-folder-line"></i>
          </div>
          <div class="dev-stat-info">
            <span class="dev-stat-value">{{ dictionaries.length }}</span>
            <span class="dev-stat-label">字典分类</span>
          </div>
        </div>
        <div class="dev-stat-card">
          <div class="dev-stat-icon items">
            <i class="ri-list-check"></i>
          </div>
          <div class="dev-stat-info">
            <span class="dev-stat-value">{{ activeDictItems.length }}</span>
            <span class="dev-stat-label">当前条目</span>
          </div>
        </div>
        <div class="dev-stat-card" v-if="activeDictId">
          <div class="dev-stat-icon active">
            <i class="ri-bookmark-line"></i>
          </div>
          <div class="dev-stat-info">
            <span class="dev-stat-value">{{ activeDictName }}</span>
            <span class="dev-stat-label">当前字典</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 内容区域 -->
    <div class="dict-layout">
      <!-- Left: Categories -->
      <div class="dict-sidebar">
        <div class="dict-sidebar-title">字典分类</div>
        <div
          v-for="dict in dictionaries"
          :key="dict.id"
          class="dict-nav-item"
          :class="{ active: activeDictId === dict.code }"
          @click="activeDictId = dict.code"
        >
          <i class="ri-folder-2-line"></i>
          <span>{{ dict.name }}</span>
        </div>
      </div>
      <!-- Right: Items -->
      <div class="dict-content">
        <div v-if="activeDictItems.length === 0" class="dict-empty">
          <div class="dict-empty-icon">
            <i class="ri-file-list-3-line"></i>
          </div>
          <h3 class="dict-empty-title">暂无数据</h3>
          <p class="dict-empty-desc">点击右上角按钮添加新条目</p>
        </div>
        <div v-else class="dict-list">
          <div v-for="item in activeDictItems" :key="item.id" class="dict-item-card">
            <div class="dict-item-icon">
              <i class="ri-price-tag-3-line"></i>
            </div>
            <div class="dict-item-info">
              <span class="dict-item-label">{{ item.label }}</span>
              <span class="dict-item-value">{{ item.value }}</span>
            </div>
            <div class="dict-item-actions">
              <button class="action-btn edit" @click="openEditModal(item)" title="编辑">
                <i class="ri-edit-line"></i>
              </button>
              <button class="action-btn delete" @click="deleteDictItem(item.id)" title="删除">
                <i class="ri-delete-bin-line"></i>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Dict Item Modal -->
    <Teleport to="body">
      <Transition name="fade">
        <div v-if="showModal" class="modal-overlay open" @click.self="showModal = false">
          <div class="modal open">
            <div class="modal-header" style="border-bottom: 1px solid var(--separator-color); padding-bottom: 16px; margin-bottom: 24px;">
              <h3 class="modal-title">{{ isEditing ? '编辑条目' : '新增条目' }}</h3>
              <button class="modal-close" @click="showModal = false">
                <i class="ri-close-line"></i>
              </button>
            </div>
            <div class="modal-body">
              <div class="form-group mb-md">
                <label class="form-label">名称 (Label)</label>
                <input
                  v-model="modalForm.label"
                  type="text"
                  class="form-input"
                  spellcheck="false"
                  autocomplete="off"
                />
              </div>
              <div class="form-group mb-md">
                <label class="form-label">值 (Value)</label>
                <input
                  v-model="modalForm.value"
                  type="text"
                  class="form-input"
                  spellcheck="false"
                  autocomplete="off"
                />
              </div>
              <div class="form-group mb-md">
                <label class="form-label">排序 (Sort)</label>
                <input
                  v-model.number="modalForm.sort"
                  type="number"
                  class="form-input"
                  spellcheck="false"
                  autocomplete="off"
                />
              </div>
            </div>
            <div class="modal-footer">
              <button class="btn btn-secondary" @click="showModal = false">取消</button>
              <button class="btn btn-primary" @click="handleModalSubmit">保存</button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<style scoped>
/* Dictionary Management Styles */
.dict-management {
  display: flex;
  flex-direction: column;
}

.dict-management .dev-header {
  padding: 1.5rem;
  border-bottom: 1px solid var(--border-color);
}

.dict-management .dev-header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.dict-management .dev-title-section {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.dict-management .dev-icon-wrapper {
  width: 3rem;
  height: 3rem;
  border-radius: 0.75rem;
  background: linear-gradient(135deg, rgba(255, 159, 10, 0.15) 0%, rgba(255, 159, 10, 0.05) 100%);
  border: 1px solid rgba(255, 159, 10, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.375rem;
  color: #FF9F0A;
}

.dict-management .dev-title-info {
  display: flex;
  flex-direction: column;
}

.dict-management .dev-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
  letter-spacing: -0.02em;
}

.dict-management .dev-subtitle {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin: 0.25rem 0 0 0;
}

.dict-management .dev-create-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.625rem 1.25rem;
  background: #FF9F0A;
  border: none;
  border-radius: 0.625rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: white;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: 0 4px 12px rgba(255, 159, 10, 0.3);
}

.dict-management .dev-create-btn:hover {
  background: #F59300;
  transform: translateY(-1px);
  box-shadow: 0 6px 16px rgba(255, 159, 10, 0.4);
}

.dict-management .dev-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1rem;
}

.dict-management .dev-stat-card {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem;
  background: var(--bg-elevated);
  border: 1px solid var(--border-color);
  border-radius: 0.75rem;
  transition: all 0.2s ease;
}

.dict-management .dev-stat-card:hover {
  border-color: rgba(255, 159, 10, 0.3);
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.dict-management .dev-stat-icon {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 0.625rem;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.125rem;
}

.dict-management .dev-stat-icon.total {
  background: linear-gradient(135deg, rgba(255, 159, 10, 0.15) 0%, rgba(255, 159, 10, 0.05) 100%);
  color: #FF9F0A;
}

.dict-management .dev-stat-icon.items {
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.15) 0%, rgba(59, 130, 246, 0.05) 100%);
  color: #3B82F6;
}

.dict-management .dev-stat-icon.active {
  background: linear-gradient(135deg, rgba(34, 197, 94, 0.15) 0%, rgba(34, 197, 94, 0.05) 100%);
  color: #22C55E;
}

.dict-management .dev-stat-info {
  display: flex;
  flex-direction: column;
}

.dict-management .dev-stat-value {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1.2;
}

.dict-management .dev-stat-value.dict-code {
  font-size: 0.875rem;
  font-weight: 600;
  font-family: monospace;
}

.dict-management .dev-stat-label {
  font-size: 0.75rem;
  color: var(--text-secondary);
  margin-top: 0.125rem;
}

.dict-layout {
  display: flex;
}

.dict-sidebar {
  width: 180px;
  border-right: 1px solid var(--border-color);
  background: transparent;
  overflow-y: auto;
  padding: 0.75rem 0;
}

.dict-sidebar-title {
  padding: 0.5rem 1rem;
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.dict-nav-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.625rem 1rem;
  cursor: pointer;
  font-size: 0.875rem;
  color: var(--text-secondary);
  border-left: 3px solid transparent;
  transition: all 0.2s;
}

.dict-nav-item i {
  font-size: 1rem;
  opacity: 0.7;
}

.dict-nav-item:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.dict-nav-item.active {
  background: rgba(255, 159, 10, 0.08);
  color: #FF9F0A;
  border-left-color: #FF9F0A;
  font-weight: 500;
}

.dict-nav-item.active i {
  opacity: 1;
}

.dict-content {
  flex: 1;
  padding: 1rem;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

.dict-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  flex: 1;
  text-align: center;
}

.dict-empty-icon {
  width: 4rem;
  height: 4rem;
  border-radius: 50%;
  background: linear-gradient(135deg, rgba(255, 159, 10, 0.1) 0%, rgba(255, 159, 10, 0.05) 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.75rem;
  color: #FF9F0A;
  margin-bottom: 1rem;
}

.dict-empty-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 0.5rem 0;
}

.dict-empty-desc {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin: 0;
}

.dict-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.dict-item-card {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.875rem 1rem;
  background: var(--bg-elevated);
  border: 1px solid var(--border-color);
  border-radius: 0.625rem;
  transition: all 0.2s;
}

.dict-item-card:hover {
  border-color: rgba(255, 159, 10, 0.3);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.06);
}

.dict-item-icon {
  width: 2.25rem;
  height: 2.25rem;
  border-radius: 0.5rem;
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.15) 0%, rgba(59, 130, 246, 0.05) 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  color: #3B82F6;
  flex-shrink: 0;
}

.dict-item-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 0.125rem;
}

.dict-item-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-primary);
}

.dict-item-value {
  font-size: 0.75rem;
  color: var(--text-secondary);
  font-family: monospace;
}

.dict-item-actions {
  display: flex;
  gap: 0.25rem;
}

.dict-item-actions .action-btn {
  width: 2rem;
  height: 2rem;
  border-radius: 0.5rem;
  border: none;
  background: transparent;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  transition: all 0.2s ease;
}

.dict-item-actions .action-btn.edit {
  color: #3B82F6;
}

.dict-item-actions .action-btn.edit:hover {
  background: rgba(59, 130, 246, 0.1);
}

.dict-item-actions .action-btn.delete {
  color: #EF4444;
}

.dict-item-actions .action-btn.delete:hover {
  background: rgba(239, 68, 68, 0.1);
}

@media (max-width: 640px) {
  .dict-management .dev-stats {
    grid-template-columns: 1fr;
  }
  
  .dict-layout {
    flex-direction: column;
  }
  
  .dict-sidebar {
    width: 100%;
    height: auto;
    border-right: none;
    border-bottom: 1px solid var(--border-color);
    display: flex;
    overflow-x: auto;
    padding: 0;
  }
  
  .dict-sidebar-title {
    display: none;
  }
  
  .dict-nav-item {
    border-left: none;
    border-bottom: 2px solid transparent;
    white-space: nowrap;
  }
  
  .dict-nav-item.active {
    border-left: none;
    border-bottom-color: #FF9F0A;
  }
}
</style>
