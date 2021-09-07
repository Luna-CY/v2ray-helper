<template>
  <img alt="404.png" src="@/assets/404.png" class="error-404">
</template>

<script lang="ts">
import {defineComponent} from 'vue'

export default defineComponent({
  name: "404"
})
</script>

<style lang="scss" scoped>
.error-404 {
  margin: 120px auto 0 auto;
  display: flex;
}

.oops {
  font-size: 32px;
  font-weight: bold;
  line-height: 40px;
  color: #1482f0;
  margin-bottom: 20px;
  animation-name: slideUp;
  animation-duration: 0.5s;
  animation-fill-mode: forwards;
}
</style>
