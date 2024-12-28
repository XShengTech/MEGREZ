<template>
  <div v-if="!editStatus" class="text-sm text-muted-color">
    <span>{{ text || props.emptyText }}</span>
    <i class="ml-1 pi pi-pencil" style="font-size: 0.75rem" @click="editStatus = true"
      v-tooltip.top="props.emptyText"></i>
  </div>
  <InputText v-else v-model="text" type="text" size="small" @blur="saveLabel" />
</template>

<script setup>
import { ref } from 'vue';

const editStatus = ref(false)

const text = defineModel()

const props = defineProps({
  emptyText: {
    type: String,
    default: 'No label'
  },

  save: {
    type: Function,
    required: true
  }
})

const saveLabel = () => {
  props.save(props.modelValue)
  editStatus.value = false
}
</script>