<template>
  <section class="mb-6 flex items-center justify-between">
    <div class="flex items-center justify-start">
      <Avatar :icon="icon" class="mr-2" size="large" shape="circle" style="background-color: #fff;" />
      <h1 class="text-3xl leading-snug font-bold -mt-[0.25rem] ml-1">
        {{ label }}
      </h1>
    </div>
    <div class="space-x-4">
      <slot />
    </div>
  </section>
</template>

<script setup>


const props = defineProps({
  label: {
    type: String,
    required: true
  },
  icon: {
    type: String,
    required: true
  }
})
</script>