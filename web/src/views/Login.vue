<template>
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="6" lg="4">
        <v-card elevation="8" class="pa-4">
          <v-card-title class="text-h4 text-center mb-4">
            <v-icon size="large" class="mr-2">mdi-shield-check</v-icon>
            Sign In
          </v-card-title>

          <v-card-subtitle class="text-center mb-4">
            Enter your organization to continue
          </v-card-subtitle>

          <v-card-text>
            <v-form @submit.prevent="handleLogin" ref="formRef">
              <v-text-field
                v-model="organizationInput"
                label="Organization Domain"
                placeholder="e.g., tesla-073704, pepsi-073704, nike-073704"
                hint="Enter your organization name (case insensitive, .localhost will be added automatically)"
                persistent-hint
                :disabled="authStore.isLoading"
                :error-messages="errorMessages"
                required
                autofocus
                prepend-inner-icon="mdi-domain"
                variant="outlined"
                class="mb-4"
                :rules="[rules.required]"
              />

              <v-alert
                v-if="authStore.error"
                type="error"
                variant="tonal"
                class="mb-4"
                closable
                @click:close="authStore.clearError()"
              >
                {{ authStore.error }}
              </v-alert>

              <v-btn
                type="submit"
                color="primary"
                block
                size="large"
                :loading="authStore.isLoading"
                :disabled="!organizationInput.trim() || authStore.isLoading"
              >
                <v-icon left class="mr-2">mdi-login</v-icon>
                Continue to Login
              </v-btn>
            </v-form>
          </v-card-text>

          <v-card-text class="text-center text-caption text-medium-emphasis mt-4">
            <v-divider class="mb-4" />
            <p>Powered by Zitadel Multi-Tenant Authentication</p>
            <p class="mt-2">
              <v-icon size="small">mdi-information</v-icon>
              You'll be redirected to your organization's login page
            </p>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
const organizationInput = ref('')
const formRef = ref()

const rules = {
  required: (value: string) => !!value || 'Organization domain is required',
}

const errorMessages = computed(() => {
  if (authStore.error) return [authStore.error]
  return []
})

async function handleLogin() {
  if (!organizationInput.value.trim()) return

  try {
    // Normalize organization domain: lowercase and ensure .localhost suffix
    let orgDomain = organizationInput.value.trim().toLowerCase()

    // Append .localhost if not already present
    if (!orgDomain.endsWith('.localhost')) {
      orgDomain = `${orgDomain}.localhost`
    }

    await authStore.login(orgDomain)
    // User will be redirected to Zitadel, no need to navigate
  } catch (error) {
    console.error('Login failed:', error)
    // Error is already set in the store
  }
}
</script>

<style scoped>
.fill-height {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}
</style>
