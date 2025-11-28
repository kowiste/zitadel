<template>
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" class="text-center">
        <v-card class="mx-auto pa-8" max-width="500" elevation="0">
          <template v-if="!error">
            <v-progress-circular
              indeterminate
              color="primary"
              size="64"
              class="mb-4"
            />
            <p class="text-h6 mb-2">Completing sign in...</p>
            <p class="text-body-2 text-medium-emphasis">Please wait while we authenticate you</p>
          </template>

          <template v-else>
            <v-icon color="error" size="64" class="mb-4">mdi-alert-circle</v-icon>
            <p class="text-h6 mb-2">Authentication Failed</p>
            <p class="text-body-2 text-medium-emphasis mb-4">{{ error }}</p>
            <v-btn color="primary" @click="handleRetry">
              <v-icon left class="mr-2">mdi-refresh</v-icon>
              Try Again
            </v-btn>
          </template>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const error = ref<string | null>(null)

onMounted(async () => {
  try {
    // Handle the OAuth callback
    await authStore.handleCallback()

    // Get the return URL from session storage or default to dashboard
    const returnUrl = sessionStorage.getItem('returnUrl') || '/dashboard'
    sessionStorage.removeItem('returnUrl')

    // Small delay for better UX
    await new Promise(resolve => setTimeout(resolve, 500))

    // Navigate to the return URL
    await router.push(returnUrl)
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'An unknown error occurred'
    console.error('Authentication callback failed:', err)
  }
})

async function handleRetry() {
  error.value = null
  await router.push({ name: 'login' })
}
</script>

<style scoped>
.fill-height {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
