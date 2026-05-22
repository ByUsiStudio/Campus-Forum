<template>
  <v-dialog
    v-model="show"
    :max-width="maxWidth"
    :persistent="persistent"
    scrollable
  >
    <v-card>
      <!-- 头部 -->
      <v-card-title class="d-flex align-center gap-3">
        <v-icon
          :color="iconColor"
          :size="iconSize"
        >
          {{ icon }}
        </v-icon>
        <span class="text-h6">{{ title }}</span>
      </v-card-title>
      
      <v-divider></v-divider>
      
      <!-- 内容 -->
      <v-card-text class="pa-4">
        <div v-if="type === 'prompt'" class="prompt-content">
          <span>{{ message }}</span>
          <v-text-field
            v-model="inputValue"
            :label="inputLabel || '输入内容'"
            :type="inputType || 'text'"
            :placeholder="inputPlaceholder"
            :rows="inputRows || 1"
            class="mt-4"
            ref="inputRef"
            @keydown.enter="handleConfirm"
          ></v-text-field>
        </div>
        <div v-else>
          <span>{{ message }}</span>
        </div>
      </v-card-text>
      
      <v-divider></v-divider>
      
      <!-- 底部按钮 -->
      <v-card-actions class="justify-end gap-3 pa-4">
        <v-btn
          v-if="showCancel"
          variant="outlined"
          color="secondary"
          @click="handleCancel"
        >
          {{ cancelText || '取消' }}
        </v-btn>
        <v-btn
          :color="confirmColor"
          :variant="confirmVariant"
          @click="handleConfirm"
        >
          {{ confirmText || '确定' }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import { ref, computed, watch, nextTick } from 'vue'

export default {
  name: 'AppModal',
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
    iconSize: {
      type: Number,
      default: 24
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
    confirmVariant: {
      type: String,
      default: 'flat'
    },
    maxWidth: {
      type: [String, Number],
      default: '500px'
    },
    persistent: {
      type: Boolean,
      default: false
    },
    // Prompt specific props
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
    
    const handleConfirm = () => {
      emit('update:show', false)
      emit('confirm', internalValue.value)
    }
    
    const handleCancel = () => {
      emit('update:show', false)
      emit('cancel')
    }
    
    return {
      inputRef,
      internalValue,
      showCancel,
      handleConfirm,
      handleCancel
    }
  }
}
</script>

<style scoped>
.prompt-content {
  display: flex;
  flex-direction: column;
}
</style>
