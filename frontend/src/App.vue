<template>
  <n-config-provider v-bind:locale="zhCN">
    <router-view/>
  </n-config-provider>
</template>

<script setup lang="ts">
import {zhCN} from "naive-ui";</script>

<style scoped>
</style>
