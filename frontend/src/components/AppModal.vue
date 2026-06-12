<script setup>
import { inject } from 'vue'
import { modalState, handleConfirm, handleCancel } from '../utils/modal'
import MarkdownViewer from './MarkdownViewer.vue'

const emit = defineEmits(['close'])

const close = () => {
  modalState.value.show = false
  emit('close')
}
</script>

<template>
  <Teleport to="body">
    <v-dialog
      v-model="modalState.show"
      max-width="500px"
      scrollable
      @close="close"
    >
      <v-card>
        <!-- 头部 -->
        <v-card-title>
          <v-icon :color="modalState.iconColor" size="24" class="mr-2">
            {{ modalState.icon }}
          </v-icon>
          {{ modalState.title }}
        </v-card-title>

        <!-- 内容 -->
        <v-card-text>
          <!-- Markdown内容 -->
          <MarkdownViewer v-if="modalState.markdown" :value="modalState.message" />
          <!-- 普通文本内容 -->
          <p v-else>{{ modalState.message }}</p>
          
          <!-- 输入框（Prompt类型） -->
          <v-text-field
            v-if="modalState.type === 'prompt'"
            v-model="modalState.inputValue"
            :label="modalState.inputLabel"
            :type="modalState.inputType"
            :placeholder="modalState.inputPlaceholder"
            :rows="modalState.inputRows"
            class="mt-4"
          ></v-text-field>
        </v-card-text>

        <!-- 底部按钮 -->
        <v-card-actions class="justify-end">
          <!-- Cancel按钮（Confirm和Prompt类型） -->
          <v-btn
            v-if="modalState.type === 'confirm' || modalState.type === 'prompt'"
            text
            color="grey"
            @click="handleCancel(); close()"
          >
            {{ modalState.cancelText }}
          </v-btn>

          <!-- Confirm按钮 -->
          <v-btn
            :color="modalState.confirmColor"
            @click="handleConfirm(modalState.type === 'prompt' ? modalState.inputValue : true); close()"
          >
            {{ modalState.confirmText }}
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </Teleport>
</template>
