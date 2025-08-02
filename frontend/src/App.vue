<script lang="ts" setup>
import GoogleDriveLayout from './components/GoogleDriveLayout.vue'
import { NConfigProvider, NMessageProvider } from 'naive-ui'</script>

<template>
  <n-config-provider>
    <n-message-provider>
      <GoogleDriveLayout/>
    </n-message-provider>
  </n-config-provider>
</template>

<style>
html, body {
  margin: 0;
  padding: 0;
  height: 100%;
  overflow: hidden;
  font-family: 'Google Sans', Roboto, Arial, sans-serif;
}

#app {
  height: 100vh;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}
</style>
