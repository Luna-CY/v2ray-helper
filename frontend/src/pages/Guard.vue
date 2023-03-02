<template>
  <div class="guard-wrap">
    <n-space vertical v-bind:wrap-item="false" style="height: 100vh" justify="center" align="center">
      <div style="width: 40vw">
        <n-input v-model:value="state.key" placeholder="口令"/>
      </div>
      <div style="width: 40vw">
        <n-button type="primary" style="width: 100%" v-on:click="verify" v-bind:loading="state.loading">验证</n-button>
      </div>
    </n-space>
  </div>
</template>

<script setup lang="ts">
import {reactive} from "vue";
import {apis} from "@/apis/main";

const state = reactive<{ key: string, loading: boolean }>({key: "", loading: false})

async function verify() {
  state.loading = true
  const response = await apis.guard({key: state.key})
  console.log(state.key)
}
</script>

<style scoped>
.guard-wrap {
  height: 100vh;
  background-color: #34495E;
}
</style>
