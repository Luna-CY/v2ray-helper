<template>
  <el-dialog :model-value="show" destroy-on-close width="80%" @close="close">
    <div>
      <el-form ref="RemoveEndpointForm" :model="form" :rules="rules" label-width="0" onsubmit="return false"
               v-on:submit="confirm">
        <el-form-item prop="password">
          <el-input v-model="form.password" autocomplete="off" :autofocus="true" placeholder="请输入删除口令"></el-input>
        </el-form-item>
      </el-form>
    </div>
    <div class="content-center">
      <el-button @click="close">取消</el-button>
      <el-button :loading="loading" type="danger" @click="confirm">确认</el-button>
    </div>
  </el-dialog>
</template>

<script lang="ts">
import {defineComponent} from "vue"
import {V2rayEndpointRemoveForm} from "@/api/v2ray_endpoint_remove"

export default defineComponent({
  name: "Remove",

  props: {
    show: Boolean,
    loading: Boolean,
  },

  emits: ['update:show', 'confirm'],

  data() {
    return {
      form: new V2rayEndpointRemoveForm(),
      rules: {
        password: [{required: true, message: '必须填写删除密码', trigger: 'blur'},],
      },
    }
  },

  methods: {
    close() {
      this.$emit('update:show', false)
      this.form.password = ""
    },

    confirm() {
      const form = this.$refs['RemoveEndpointForm'] as any

      form.validate((valid: boolean) => {
        if (!valid) {
          return
        }

        this.$emit('confirm', this.form.password)
      })
    },
  },
})
</script>

<style scoped>

</style>
