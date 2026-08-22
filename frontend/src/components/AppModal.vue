<template>
  <el-dialog
    :model-value="show"
    :width="'500px'"
    :close-on-click-modal="false"
    append-to-body
    @update:model-value="$emit('update:show', $event)"
    @closed="resetInternal"
  >
    <template #header>
      <div class="modal-title">
        <el-icon :color="iconColor" :size="iconSize" class="title-icon">
          <WarningFilled v-if="icon === 'mdi-alert-circle' || icon === 'warning'" />
          <InfoFilled v-else-if="icon === 'mdi-information' || icon === 'info'" />
          <CircleCheckFilled v-else-if="icon === 'mdi-check-circle' || icon === 'success'" />
          <EditPen v-else-if="icon === 'mdi-edit' || icon === 'edit'" />
          <WarningFilled v-else />
        </el-icon>
        <span>{{ title }}</span>
      </div>
    </template>

    <!-- 内容 -->
    <div class="modal-body">
      <div v-if="type === 'prompt'" class="prompt-content">
        <span class="prompt-message">{{ message }}</span>
        <el-input
          v-model="internalValue"
          :label="inputLabel || '输入内容'"
          :type="inputType === 'textarea' ? 'textarea' : 'text'"
          :placeholder="inputPlaceholder"
          :rows="inputRows || 1"
          ref="inputRef"
          class="mt-4"
          @keydown.enter="handleConfirm"
        ></el-input>
      </div>
      <div v-else>
        <span>{{ message }}</span>
      </div>
    </div>

    <!-- 底部按钮 -->
    <template #footer>
      <div class="modal-actions">
        <el-button v-if="showCancel" @click="handleCancel">
          {{ cancelText || '取消' }}
        </el-button>
        <el-button :type="confirmButtonType" @click="handleConfirm">
          {{ confirmText || '确定' }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script>
import { ref, computed, watch, nextTick } from 'vue'
import {
  WarningFilled,
  InfoFilled,
  CircleCheckFilled,
  EditPen
} from '@element-plus/icons-vue'

export default {
  name: 'AppModal',
  components: {
    WarningFilled,
    InfoFilled,
    CircleCheckFilled,
    EditPen
  },
  props: {
    show: {
      type: Boolean,
      default: false
    },
    type: {
      type: String,
      default: 'alert', // alert, confirm, prompt
      validator: (value) => ['alert', 'confirm', 'prompt'].includes(value)
    },
    title: {
      type: String,
      default: ''
    },
    message: {
      type: String,
      default: ''
    },
    icon: {
      type: String,
      default: 'mdi-alert-circle'
    },
    iconColor: {
      type: String,
      default: 'warning'
    },
    confirmText: {
      type: String,
      default: '确定'
    },
    cancelText: {
      type: String,
      default: '取消'
    },
    confirmColor: {
      type: String,
      default: 'primary'
    },
    inputValue: {
      type: String,
      default: ''
    },
    inputLabel: {
      type: String,
      default: ''
    },
    inputType: {
      type: String,
      default: 'text'
    },
    inputPlaceholder: {
      type: String,
      default: ''
    },
    inputRows: {
      type: Number,
      default: 1
    }
  },
  emits: ['update:show', 'confirm', 'cancel'],
  setup(props, { emit }) {
    const inputRef = ref(null)
    const internalValue = ref(props.inputValue)

    const showCancel = computed(() => {
      return props.type === 'confirm' || props.type === 'prompt'
    })

    const confirmButtonType = computed(() => {
      const color = props.confirmColor
      if (color === 'error' || color === 'danger') return 'danger'
      if (color === 'success') return 'success'
      if (color === 'warning') return 'warning'
      if (color === 'secondary') return 'info'
      return 'primary'
    })

    watch(() => props.inputValue, (val) => {
      internalValue.value = val
    })

    watch(() => props.show, (val) => {
      if (val && props.type === 'prompt') {
        nextTick(() => {
          inputRef.value?.focus()
        })
      }
    })

    const resetInternal = () => {
      internalValue.value = props.inputValue
    }

    const handleConfirm = () => {
      emit('confirm', internalValue.value)
      emit('update:show', false)
    }

    const handleCancel = () => {
      emit('cancel')
      emit('update:show', false)
    }

    return {
      inputRef,
      internalValue,
      showCancel,
      confirmButtonType,
      resetInternal,
      handleConfirm,
      handleCancel
    }
  }
}
</script>

<style scoped>
.modal-title {
  display: flex;
  align-items: center;
  gap: 12px;
  font-weight: 600;
  font-size: 16px;
  color: #1a1a2e;
}

.title-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.prompt-content {
  display: flex;
  flex-direction: column;
}

.prompt-message {
  margin-bottom: 12px;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
