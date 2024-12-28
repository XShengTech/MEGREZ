<script setup>
import { useLayout } from '@/layout/composables/layout';

const { toggleDarkMode, isDarkTheme } = useLayout();
</script>

<template>
  <div class="fixed flex gap-4 top-8 right-8">
    <Button type="button" @click="toggleDarkMode" rounded :icon="isDarkTheme ? 'pi pi-moon' : 'pi pi-sun'"
      severity="secondary" />
  </div>
</template>
