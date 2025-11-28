<template>
  <v-app>
    <AppHeader v-if="showHeader" />

    <v-main>
      <slot />
    </v-main>
  </v-app>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import AppHeader from './AppHeader.vue'

const route = useRoute()

// Don't show header on login and callback pages
const showHeader = computed(() => {
  return route.name !== 'login' && route.name !== 'auth-callback'
})
</script>
