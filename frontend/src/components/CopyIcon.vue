<template>
  <i v-if="!clickStatus" class="pi pi-copy" @click="copyToClipboard(text)" />
  <i v-else class="pi pi-check text-green-500" @click="copyToClipboard(text)" />
</template>

<script setup>
import { computed, ref } from 'vue';

const props = defineProps({
  text: {
    type: String,
    required: true
  },
})

const clickStatus = ref(false)
const timeoutId = ref(null)

const spanClass = computed(() => `inline-flex justify-center items-center ${props.w} ${props.h}`)
const spanClassCheck = computed(() => `inline-flex justify-center items-center ${props.w} ${props.h} text-green-500`)

const iconSize = computed(() => props.size ?? 16)

const copyToClipboard = (text) => {
  clearTimeout(timeoutId.value)
  if (navigator.clipboard) {
    navigator.clipboard.writeText(text).then(() => {
      console.log('Text copied to clipboard')
    }).catch(err => {
      console.error('Failed to copy text: ', err)
    });
  } else {
    const textarea = document.createElement('textarea')
    textarea.value = text
    document.body.appendChild(textarea)
    textarea.select();
    try {
      document.execCommand('copy')
      console.log('Text copied to clipboard')
    } catch (err) {
      console.error('Failed to copy text: ', err)
    }
    document.body.removeChild(textarea)
  }
  clickStatus.value = true
  timeoutId.vaue = setTimeout(() => {
    clickStatus.value = false
  }, 1000)
}
</script>