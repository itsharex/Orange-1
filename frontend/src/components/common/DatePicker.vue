<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import dayjs from 'dayjs'

const props = defineProps<{
  modelValue: string
  placeholder?: string
  required?: boolean
  disabled?: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

const isOpen = ref(false)
const wrapperRef = ref<HTMLElement | null>(null)

// Current view date (for the calendar grid)
const viewDate = ref(dayjs())

// Initialize viewDate from modelValue if present
watch(() => props.modelValue, (newVal) => {
  if (newVal) {
    const d = dayjs(newVal)
    if (d.isValid()) {
      viewDate.value = d
    }
  } else {
      // If cleared, maybe reset viewDate to today? Or keep current?
      // Keeping current viewDate is usually better UX.
      if (!viewDate.value.isValid()) viewDate.value = dayjs()
  }
}, { immediate: true })

const formattedValue = computed(() => {
  if (!props.modelValue) return ''
  // Ensure we display correctly
  return dayjs(props.modelValue).format('YYYY-MM-DD')
})

const year = computed(() => viewDate.value.year())
const month = computed(() => viewDate.value.month()) // 0-indexed

// Calendar Grid Logic
const days = computed(() => {
  const startOfMonth = viewDate.value.startOf('month')
  const startDayOfWeek = startOfMonth.day() // 0 (Sunday) to 6 (Saturday)
  
  const daysArray = []
  
  // Previous month days padding
  const prevMonth = viewDate.value.subtract(1, 'month')
  const daysInPrevMonth = prevMonth.daysInMonth()
  for (let i = startDayOfWeek - 1; i >= 0; i--) {
    daysArray.push({
      date: prevMonth.date(daysInPrevMonth - i),
      isCurrentMonth: false,
      isToday: false,
      isSelected: false
    })
  }
  
  // Current month days
  const daysInMonth = viewDate.value.daysInMonth()
  const today = dayjs()
  for (let i = 1; i <= daysInMonth; i++) {
    const d = viewDate.value.date(i)
    daysArray.push({
      date: d,
      isCurrentMonth: true,
      isToday: d.isSame(today, 'day'),
      isSelected: props.modelValue ? d.isSame(dayjs(props.modelValue), 'day') : false
    })
  }
  
  // Next month days padding (to fill 6 rows = 42 cells)
  const remaining = 42 - daysArray.length
  const nextMonth = viewDate.value.add(1, 'month')
  for (let i = 1; i <= remaining; i++) {
     daysArray.push({
      date: nextMonth.date(i),
      isCurrentMonth: false,
      isToday: false,
      isSelected: false
    })
  }
  
  return daysArray
})

const weekDays = ['日', '一', '二', '三', '四', '五', '六']

// Actions
const togglePicker = () => {
    if (props.disabled) return
    isOpen.value = !isOpen.value
}

const selectDate = (date: dayjs.Dayjs) => {
    emit('update:modelValue', date.format('YYYY-MM-DD'))
    isOpen.value = false
    // Also update view to this date?
    viewDate.value = date
}

const changeMonth = (delta: number) => {
    viewDate.value = viewDate.value.add(delta, 'month')
}

const changeYear = (delta: number) => {
    viewDate.value = viewDate.value.add(delta, 'year')
}

// Click Outside
const handleClickOutside = (event: MouseEvent) => {
  if (wrapperRef.value && !wrapperRef.value.contains(event.target as Node)) {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<template>
  <div class="date-picker-wrapper" ref="wrapperRef">
    <!-- Input Trigger -->
    <div class="input-trigger" @click="togglePicker" :class="{ 'is-disabled': disabled }">
        <input 
            type="text" 
            readonly 
            :value="formattedValue" 
            :placeholder="placeholder || '请选择日期'"
            :required="required"
            class="readonly-input"
        />
        <i class="ri-calendar-line icon"></i>
    </div>

    <!-- Dropdown Panel -->
    <transition name="fade-slide">
        <div v-if="isOpen" class="picker-panel glass-panel">
            <!-- Header -->
            <div class="picker-header">
                <button type="button" class="nav-btn" @click.stop="changeYear(-1)"><i class="ri-arrow-left-double-line"></i></button>
                <button type="button" class="nav-btn" @click.stop="changeMonth(-1)"><i class="ri-arrow-left-s-line"></i></button>
                
                <div class="current-date">
                    <span class="year">{{ year }}年</span>
                    <span class="month">{{ month + 1 }}月</span>
                </div>
                
                <button type="button" class="nav-btn" @click.stop="changeMonth(1)"><i class="ri-arrow-right-s-line"></i></button>
                <button type="button" class="nav-btn" @click.stop="changeYear(1)"><i class="ri-arrow-right-double-line"></i></button>
            </div>
            
            <!-- Weekdays -->
            <div class="weekdays-row">
                <span v-for="day in weekDays" :key="day" class="weekday">{{ day }}</span>
            </div>
            
            <!-- Days Grid -->
            <div class="days-grid">
                <div 
                    v-for="(item, idx) in days" 
                    :key="idx" 
                    class="day-cell"
                    :class="{ 
                        'is-current-month': item.isCurrentMonth,
                        'is-prev-next': !item.isCurrentMonth,
                        'is-today': item.isToday,
                        'is-selected': item.isSelected
                    }"
                    @click.stop="selectDate(item.date)"
                >
                    {{ item.date.date() }}
                </div>
            </div>
        </div>
    </transition>
  </div>
</template>

<style scoped>
.date-picker-wrapper {
    position: relative;
    width: 100%;
}

.input-trigger {
    position: relative;
    cursor: pointer;
}

.readonly-input {
    width: 100%;
    padding: 10px 14px;
    padding-right: 36px;
    font-size: 14px;
    color: var(--text-primary);
    background: var(--bg-base);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-sm);
    outline: none;
    cursor: pointer; /* Important for readonly feel */
    transition: all 0.2s;
}

.readonly-input:focus {
    border-color: var(--color-primary);
    box-shadow: 0 0 0 3px rgba(255, 159, 10, 0.1);
}

.icon {
    position: absolute;
    right: 12px;
    top: 50%;
    transform: translateY(-50%);
    color: var(--text-secondary);
    pointer-events: none;
    font-size: 16px;
}

.picker-panel {
    position: absolute;
    top: calc(100% + 8px);
    left: 0;
    width: 280px;
    z-index: 100;
    padding: 16px;
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(20px);
    border: 1px solid rgba(255, 255, 255, 0.3);
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
    border-radius: 12px;
}

[data-theme='dark'] .picker-panel {
    background: rgba(30, 30, 30, 0.95);
    border-color: rgba(255, 255, 255, 0.05);
}

/* Header */
.picker-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
}

.current-date {
    font-weight: 600;
    display: flex;
    gap: 4px;
    color: var(--text-primary);
}

.nav-btn {
    width: 28px;
    height: 28px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 6px;
    color: var(--text-secondary);
    transition: all 0.2s;
}

.nav-btn:hover {
    background: rgba(0,0,0,0.05);
    color: var(--color-primary);
}

[data-theme='dark'] .nav-btn:hover {
    background: rgba(255,255,255,0.1);
}

/* Weekdays */
.weekdays-row {
    display: grid;
    grid-template-columns: repeat(7, 1fr);
    margin-bottom: 8px;
}

.weekday {
    text-align: center;
    font-size: 12px;
    color: var(--text-secondary);
    font-weight: 500;
}

/* Days Grid */
.days-grid {
    display: grid;
    grid-template-columns: repeat(7, 1fr);
    row-gap: 4px;
}

.day-cell {
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 13px;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s;
    color: var(--text-primary);
}

.day-cell.is-prev-next {
    color: var(--text-placeholder);
    opacity: 0.5;
}

.day-cell:hover {
    background: rgba(var(--color-primary-rgb), 0.1);
    color: var(--color-primary);
}

.day-cell.is-today {
    color: var(--color-primary);
    font-weight: 600;
}

.day-cell.is-selected {
    background: var(--color-primary);
    color: white;
    font-weight: 600;
    box-shadow: 0 2px 8px rgba(var(--color-primary-rgb), 0.4);
}

.day-cell.is-selected:hover {
    background: var(--color-primary);
    opacity: 0.9;
}

/* Animations */
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

.fade-slide-enter-from,
.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-8px) scale(0.98);
}
</style>
