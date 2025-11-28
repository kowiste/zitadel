<template>
  <v-container class="py-8">
    <v-row>
      <v-col cols="12">
        <h1 class="text-h3 mb-2">
          <v-icon size="large" class="mr-2">mdi-view-dashboard</v-icon>
          Dashboard
        </h1>
        <p class="text-h6 text-medium-emphasis mb-6">
          Welcome, {{ authStore.userProfile?.name || 'User' }}!
        </p>
      </v-col>
    </v-row>

    <v-row>
      <!-- User Info Card -->
      <v-col cols="12" md="6">
        <v-card elevation="2">
          <v-card-title class="bg-primary">
            <v-icon class="mr-2">mdi-account-circle</v-icon>
            User Information
          </v-card-title>
          <v-card-text class="pt-4">
            <v-list>
              <v-list-item>
                <template v-slot:prepend>
                  <v-icon>mdi-account</v-icon>
                </template>
                <v-list-item-title>Name</v-list-item-title>
                <v-list-item-subtitle class="text-wrap">
                  {{ authStore.userProfile?.name || 'N/A' }}
                </v-list-item-subtitle>
              </v-list-item>

              <v-list-item>
                <template v-slot:prepend>
                  <v-icon>mdi-email</v-icon>
                </template>
                <v-list-item-title>Email</v-list-item-title>
                <v-list-item-subtitle class="text-wrap">
                  {{ authStore.userProfile?.email || 'N/A' }}
                </v-list-item-subtitle>
              </v-list-item>

              <v-list-item>
                <template v-slot:prepend>
                  <v-icon>mdi-identifier</v-icon>
                </template>
                <v-list-item-title>User ID</v-list-item-title>
                <v-list-item-subtitle class="text-wrap font-mono">
                  {{ authStore.userProfile?.sub || 'N/A' }}
                </v-list-item-subtitle>
              </v-list-item>

              <v-list-item v-if="authStore.userProfile?.email_verified !== undefined">
                <template v-slot:prepend>
                  <v-icon>{{ authStore.userProfile?.email_verified ? 'mdi-check-circle' : 'mdi-alert-circle' }}</v-icon>
                </template>
                <v-list-item-title>Email Verified</v-list-item-title>
                <v-list-item-subtitle>
                  <v-chip :color="authStore.userProfile?.email_verified ? 'success' : 'warning'" size="small">
                    {{ authStore.userProfile?.email_verified ? 'Yes' : 'No' }}
                  </v-chip>
                </v-list-item-subtitle>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Organization Info Card -->
      <v-col cols="12" md="6">
        <v-card elevation="2">
          <v-card-title class="bg-primary">
            <v-icon class="mr-2">mdi-domain</v-icon>
            Organization Details
          </v-card-title>
          <v-card-text class="pt-4">
            <v-list>
              <v-list-item>
                <template v-slot:prepend>
                  <v-icon>mdi-domain</v-icon>
                </template>
                <v-list-item-title>Organization Domain</v-list-item-title>
                <v-list-item-subtitle class="text-wrap">
                  {{ authStore.selectedOrganization }}
                </v-list-item-subtitle>
              </v-list-item>

              <v-list-item>
                <template v-slot:prepend>
                  <v-icon>mdi-clock-check</v-icon>
                </template>
                <v-list-item-title>Authenticated</v-list-item-title>
                <v-list-item-subtitle>
                  <v-chip color="success" size="small">
                    <v-icon start size="small">mdi-check</v-icon>
                    Active
                  </v-chip>
                </v-list-item-subtitle>
              </v-list-item>

              <v-list-item v-if="tokenExpiresAt">
                <template v-slot:prepend>
                  <v-icon>mdi-clock-outline</v-icon>
                </template>
                <v-list-item-title>Session Expires</v-list-item-title>
                <v-list-item-subtitle>
                  {{ tokenExpiresAt }}
                </v-list-item-subtitle>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12">
        <v-card elevation="2">
          <v-card-title class="bg-primary">
            <v-icon class="mr-2">mdi-information</v-icon>
            Quick Actions
          </v-card-title>
          <v-card-text class="pt-4">
            <v-row>
              <v-col cols="12" sm="6" md="4">
                <v-btn
                  :to="{ name: 'profile' }"
                  color="primary"
                  block
                  size="large"
                  variant="outlined"
                >
                  <v-icon left class="mr-2">mdi-account</v-icon>
                  View Profile
                </v-btn>
              </v-col>
            </v-row>
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

const tokenExpiresAt = computed(() => {
  if (!authStore.user?.expires_at) return null
  const date = new Date(authStore.user.expires_at * 1000)
  return date.toLocaleString()
})
</script>

<style scoped>
.font-mono {
  font-family: monospace;
  font-size: 0.875rem;
}
</style>
