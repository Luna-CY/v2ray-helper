<template>
  <div class="login-box">
    <div class="login-form">
      <el-form ref="LoginForm" :model="form" :rules="rules" label-width="0" onsubmit="return false" v-on:submit="login">
        <el-form-item prop="key">
          <el-input v-model="form.key" autocomplete="off" placeholder="口令" prefix-icon="el-icon-key"
                    size="medium">
            <template #append>
              <el-button :loading="logging" icon="el-icon-check" type="primary" @click="login"></el-button>
            </template>
          </el-input>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script lang="ts">
import {defineComponent} from 'vue'
import axios, {AxiosResponse} from "axios"
import {API_LOGIN, LoginForm, LoginResponse} from "@/api/login"
import {StoryStateToken} from "@/store"

const md5 = require('md5')

export default defineComponent({
  name: "Login",

  data() {
    return {
      form: new LoginForm(),
      rules: {
        key: [{required: true, message: '请输入口令', trigger: 'blur'},],
      },
      logging: false,
    }
  },

  methods: {
    login() {
      const form = this.$refs['LoginForm'] as any

      form.validate((valid: boolean) => {
        if (!valid) {
          return false
        }

        const body = new LoginForm()
        body.key = md5(this.form.key)

        this.logging = true
        axios.post(API_LOGIN, body).then((response: AxiosResponse<LoginResponse>) => {
          this.logging = false

          if (0 != response.data.code) {
            this.$message.error(response.data.message)

            return
          }

          const token = new StoryStateToken()
          token.token = response.data.data.token
          token.expired = response.data.data.expired

          this.$store.commit('token', token)
          const redirect = this.$route.query.redirect as string

          window.location.href = redirect ? redirect : '/'
        })
      })

      return false
    }
  },
})
</script>

<style lang="scss">
body {
  background-color: #2d3a4b;
}
</style>

<style lang="scss" scoped>
.login-box {
  color: #eee;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 60vh;

  .login-form {
    min-width: 20%;

    .login-form-title {
      font-size: 2em;
      margin-bottom: 45px;
      text-align: center;
    }

    .el-form-item {
      display: flex;

      ::v-deep .el-form-item__content {
        flex-grow: 1;

        .el-input__inner {
          color: #eee;
          background-color: #2d3a4b;
          border: 1px solid hsla(0, 0%, 100%, .1);
        }
      }
    }
  }
}
</style>
