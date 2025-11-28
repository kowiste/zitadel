<template>
  <v-app-bar color="primary" prominent>
    <v-app-bar-title>
      <v-icon class="mr-2">mdi-shield-check</v-icon>
      Zitadel Multi-Tenant App
    </v-app-bar-title>

    <v-spacer />

    <template v-if="authStore.isAuthenticated">
      <v-chip class="mr-4" prepend-icon="mdi-domain">
        {{ authStore.selectedOrganization }}
      </v-chip>

      <v-menu>
        <template v-slot:activator="{ props }">
          <v-btn icon v-bind="props">
            <v-icon>mdi-account-circle</v-icon>
          </v-btn>
        </template>

        <v-list>
          <v-list-item>
            <v-list-item-title>{{ authStore.userProfile?.name || 'User' }}</v-list-item-title>
            <v-list-item-subtitle>{{ authStore.userProfile?.email }}</v-list-item-subtitle>
          </v-list-item>

          <v-divider />

          <v-list-item :to="{ name: 'dashboard' }" prepend-icon="mdi-view-dashboard">
            <v-list-item-title>Dashboard</v-list-item-title>
          </v-list-item>

          <v-list-item :to="{ name: 'profile' }" prepend-icon="mdi-account">
            <v-list-item-title>Profile</v-list-item-title>
          </v-list-item>

          <v-divider />

          <v-list-item @click="handleLogout" prepend-icon="mdi-logout">
            <v-list-item-title>Logout</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
    </template>
  </v-app-bar>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const router = useRouter()

async function handleLogout() {
  await authStore.logout()
  await router.push({ name: 'login' })
}
</script>
