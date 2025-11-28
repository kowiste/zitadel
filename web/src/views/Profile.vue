<template>
  <v-container class="py-8">
    <v-row>
      <v-col cols="12">
        <h1 class="text-h3 mb-2">
          <v-icon size="large" class="mr-2">mdi-account</v-icon>
          Profile
        </h1>
        <p class="text-body-1 text-medium-emphasis mb-6">
          Your account information and settings
        </p>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12" md="8">
        <!-- Profile Information -->
        <v-card elevation="2" class="mb-4">
          <v-card-title class="bg-primary">
            <v-icon class="mr-2">mdi-account-details</v-icon>
            Profile Information
          </v-card-title>
          <v-card-text class="pt-4">
            <v-list>
              <v-list-item v-if="authStore.userProfile?.given_name">
                <template v-slot:prepend>
                  <v-avatar color="primary">
                    <v-icon>mdi-account</v-icon>
                  </v-avatar>
                </template>
                <v-list-item-title>Given Name</v-list-item-title>
                <v-list-item-subtitle>{{ authStore.userProfile.given_name }}</v-list-item-subtitle>
              </v-list-item>

              <v-list-item v-if="authStore.userProfile?.family_name">
                <template v-slot:prepend>
                  <v-avatar color="primary">
                    <v-icon>mdi-account</v-icon>
                  </v-avatar>
                </template>
                <v-list-item-title>Family Name</v-list-item-title>
                <v-list-item-subtitle>{{ authStore.userProfile.family_name }}</v-list-item-subtitle>
              </v-list-item>

              <v-list-item>
                <template v-slot:prepend>
                  <v-avatar color="primary">
                    <v-icon>mdi-email</v-icon>
                  </v-avatar>
                </template>
                <v-list-item-title>Email Address</v-list-item-title>
                <v-list-item-subtitle>{{ authStore.userProfile?.email || 'N/A' }}</v-list-item-subtitle>
              </v-list-item>

              <v-list-item v-if="authStore.userProfile?.preferred_username">
                <template v-slot:prepend>
                  <v-avatar color="primary">
                    <v-icon>mdi-at</v-icon>
                  </v-avatar>
                </template>
                <v-list-item-title>Username</v-list-item-title>
                <v-list-item-subtitle>{{ authStore.userProfile.preferred_username }}</v-list-item-subtitle>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>

        <!-- Token Information -->
        <v-card elevation="2">
          <v-card-title class="bg-primary">
            <v-icon class="mr-2">mdi-key</v-icon>
            Session Information
          </v-card-title>
          <v-card-text class="pt-4">
            <v-list>
              <v-list-item>
                <template v-slot:prepend>
                  <v-icon>mdi-identifier</v-icon>
                </template>
                <v-list-item-title>Subject ID</v-list-item-title>
                <v-list-item-subtitle class="font-mono">
                  {{ authStore.userProfile?.sub || 'N/A' }}
                </v-list-item-subtitle>
              </v-list-item>

              <v-list-item>
                <template v-slot:prepend>
                  <v-icon>mdi-domain</v-icon>
                </template>
                <v-list-item-title>Organization</v-list-item-title>
                <v-list-item-subtitle>
                  {{ authStore.selectedOrganization }}
                </v-list-item-subtitle>
              </v-list-item>

              <v-list-item v-if="tokenIssuedAt">
                <template v-slot:prepend>
                  <v-icon>mdi-clock-start</v-icon>
                </template>
                <v-list-item-title>Token Issued At</v-list-item-title>
                <v-list-item-subtitle>
                  {{ tokenIssuedAt }}
                </v-list-item-subtitle>
              </v-list-item>

              <v-list-item v-if="tokenExpiresAt">
                <template v-slot:prepend>
                  <v-icon>mdi-clock-end</v-icon>
                </template>
                <v-list-item-title>Token Expires At</v-list-item-title>
                <v-list-item-subtitle>
                  {{ tokenExpiresAt }}
                </v-list-item-subtitle>
              </v-list-item>

              <v-list-item>
                <template v-slot:prepend>
                  <v-icon :color="tokenStatus.color">{{ tokenStatus.icon }}</v-icon>
                </template>
                <v-list-item-title>Token Status</v-list-item-title>
                <v-list-item-subtitle>
                  <v-chip :color="tokenStatus.color" size="small">
                    {{ tokenStatus.text }}
                  </v-chip>
                </v-list-item-subtitle>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" md="4">
        <!-- Profile Actions -->
        <v-card elevation="2" class="mb-4">
          <v-card-title class="bg-primary">
            <v-icon class="mr-2">mdi-cog</v-icon>
            Actions
          </v-card-title>
          <v-card-text class="pt-4">
            <v-btn
              :to="{ name: 'dashboard' }"
              color="primary"
              block
              variant="outlined"
              class="mb-2"
            >
              <v-icon left class="mr-2">mdi-view-dashboard</v-icon>
              Back to Dashboard
            </v-btn>
          </v-card-text>
        </v-card>

        <!-- Additional Profile Info -->
        <v-card elevation="2">
          <v-card-title class="bg-primary">
            <v-icon class="mr-2">mdi-information</v-icon>
            Additional Claims
          </v-card-title>
          <v-card-text class="pt-4">
            <v-expansion-panels>
              <v-expansion-panel>
                <v-expansion-panel-title>
                  View All Token Claims
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <pre class="text-caption">{{ JSON.stringify(authStore.userProfile, null, 2) }}</pre>
                </v-expansion-panel-text>
              </v-expansion-panel>
            </v-expansion-panels>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const tokenIssuedAt = computed(() => {
  if (!authStore.user?.profile?.iat) return null
  const date = new Date(authStore.user.profile.iat * 1000)
  return date.toLocaleString()
})

const tokenExpiresAt = computed(() => {
  if (!authStore.user?.expires_at) return null
  const date = new Date(authStore.user.expires_at * 1000)
  return date.toLocaleString()
})

const tokenStatus = computed(() => {
  if (!authStore.user) {
    return { text: 'No Token', color: 'grey', icon: 'mdi-close-circle' }
  }
  if (authStore.user.expired) {
    return { text: 'Expired', color: 'error', icon: 'mdi-alert-circle' }
  }
  return { text: 'Valid', color: 'success', icon: 'mdi-check-circle' }
})
</script>

<style scoped>
.font-mono {
  font-family: monospace;
  font-size: 0.875rem;
}

pre {
  overflow-x: auto;
  background-color: #f5f5f5;
  padding: 1rem;
  border-radius: 4px;
  max-height: 300px;
  overflow-y: auto;
}
</style>
