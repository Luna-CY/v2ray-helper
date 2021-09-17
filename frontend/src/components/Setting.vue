<template>
  <el-dialog :model-value="show" destroy-on-close width="80%" @close="close">
    <el-form onsubmit="return false" label-width="120px">
      <el-form-item label="监听端口">
        <div class="inline-button">
          <el-input v-model="listen" placeholder="V2rayHelper监听的端口，1-65535的整数，不能使用80与443端口"></el-input>
          <el-button type="primary" class="margin-left" @click="save('listen', listen)">保存</el-button>
        </div>
      </el-form-item>
      <div class="inline-form-item-2">
        <el-form-item class="form-item-0" label="HTTPS">
          <el-switch v-model="enable_https"></el-switch>
        </el-form-item>
        <el-form-item class="form-item-1" label="域名" label-width="60px">
          <div class="inline-button">
            <el-input v-model="https_host" :disabled="!enable_https"
                      placeholder="HTTPS域名，将使用该域名申请证书，确保该域名已解析到当前服务器"></el-input>
            <el-button type="primary" class="margin-left" @click="save('https-host', https_host)">保存</el-button>
          </div>
        </el-form-item>
      </div>
      <el-form-item label="邮件地址">
        <div class="inline-button">
          <el-input v-model="email" placeholder="用于注册HTTPS证书，默认地址为: myself@v2ray-helper.net"></el-input>
          <el-button type="primary" class="margin-left" @click="save('email', email)">保存</el-button>
        </div>
      </el-form-item>
    </el-form>
    <div class="content-center">
      <el-button @click="close">关闭</el-button>
    </div>
  </el-dialog>
</template>

<script lang="ts">
import {defineComponent} from "vue"
import {API_SAVE_META_INFO, SaveMetaInfoForm} from "@/api/save_meta_info";
import axios, {AxiosResponse} from "axios";
import {BaseResponse} from "@/api/base";

export default defineComponent({
  name: "Setting",

  props: {
    show: Boolean,
  },

  emits: ['update:show'],

  watch: {
    show() {
      this.listen = this.$store.getters.local.listen
      this.enable_https = this.$store.getters.local.enable_https
      this.https_host = this.$store.getters.local.https_host
      this.email = this.$store.getters.local.email
    },

    enable_https() {
      if (!this.enable_https) {
        this.https_host = ""
      }
    }
  },

  data() {
    return {
      listen: 8888,
      enable_https: false,
      https_host: "",
      email: "",
    }
  },

  methods: {
    close() {
      this.$emit('update:show', false)
    },

    save(key: string, value: any) {
      switch (key) {
        case 'listen':
          let port = parseInt(value.toString())
          if (1 > port || 65535 < port || 80 == port || 443 == port) {
            this.$message.error('无效的端口号')

            return
          }

          break
        case 'https-host':
          value = value.trim()
          if (this.enable_https && "" == value.trim()) {
            this.$message.error('启用HTTPS时必须填写域名')
          }

          break
        case 'email':
          value = value.trim()

          break
        default:
          return
      }

      let form = new SaveMetaInfoForm()
      form.key = key
      form.value = value

      axios.post(API_SAVE_META_INFO, form).then((response: AxiosResponse<BaseResponse>) => {
        if (0 != response.data.code) {
          this.$message.error(response.data.message)

          return
        }

        let story = this.$store.getters.local
        story.listen = parseInt(this.listen.toString())
        story.enable_https = this.enable_https
        story.https_host = this.https_host
        story.email = this.email

        this.$store.commit('local', story)

        this.$message.success('保存成功，重启服务后生效')
      })
    }
  },
})
</script>

<style scoped lang="scss">
.inline-form-item-2 {
  display: flex;

  .form-item {
    width: 50%;
  }

  .form-item-0 {
    flex-grow: 0;
  }

  .form-item-1 {
    flex-grow: 1;
  }
}

.inline-button {
  display: flex;

  .el-input {
    flex-grow: 1;
  }
}
</style>
