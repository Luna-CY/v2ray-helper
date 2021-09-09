<template>
  <el-dialog width="60%" :model-value="show" @close="close" :close-on-click-modal="false"
             :close-on-press-escape="false" destroy-on-close>
    <el-form :model="form" label-width="120px" ref="NewV2rayEndpointForm" :rules="rules">
      <el-form-item label="服务商" prop="cloud">
        <el-select v-model="form.cloud">
          <el-option :value="0" label="请选择服务商"></el-option>
          <el-option :value="1" label="Vultr"></el-option>
          <el-option :value="2" label="阿里云"></el-option>
          <el-option :value="3" label="腾讯云"></el-option>
          <el-option :value="4" label="华为云"></el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="节点位置" prop="endpoint">
        <el-select v-model="form.endpoint">
          <el-option :value="0" label="请选择节点位置"></el-option>
          <el-option :value="1" label="日本"></el-option>
          <el-option :value="2" label="香港"></el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="带宽上限">
        <el-input v-model="form.rate" placeholder="e.g.: 100M"></el-input>
      </el-form-item>
      <el-form-item label="备注信息">
        <el-input v-model="form.remark" placeholder="备注信息"></el-input>
      </el-form-item>
      <el-form-item label="主机HOST" prop="host">
        <el-input v-model="form.host" placeholder="主机地址: IP或域名"></el-input>
      </el-form-item>
      <el-form-item label="端口" prop="port">
        <el-input v-model="form.port" placeholder="端口"></el-input>
      </el-form-item>
      <el-form-item label="用户ID" prop="user_id">
        <el-input v-model="form.user_id" placeholder="用户身份ID"></el-input>
      </el-form-item>
      <el-form-item label="AlterId" prop="alter_id">
        <el-input v-model="form.alter_id" placeholder="e.g.: 64"></el-input>
      </el-form-item>
      <el-form-item label="SNI">
        <el-input v-model="form.sni" placeholder=""></el-input>
      </el-form-item>
      <el-form-item label="使用TLS">
        <el-switch v-model="form.use_tls"></el-switch>
      </el-form-item>
      <el-form-item label="传输方式" prop="transport_type">
        <el-select v-model="form.transport_type">
          <el-option :value="1" label="TCP"></el-option>
          <el-option :value="2" label="WebSocket"></el-option>
          <el-option :value="3" label="KCP"></el-option>
          <el-option :value="4" label="HTTP2"></el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="伪装类型" v-if="1 === parseInt(form.transport_type.toString())">
        <el-select v-model="form.tcp.type">
          <el-option value="none" label="none"></el-option>
          <el-option value="http" label="http"></el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="伪装类型" v-if="3 === parseInt(form.transport_type.toString())">
        <el-select v-model="form.kcp.type">
          <el-option value="none" label="none"></el-option>
          <el-option value="srtp" label="srtp"></el-option>
          <el-option value="utp" label="utp"></el-option>
          <el-option value="wechat-video" label="wechat-video"></el-option>
          <el-option value="dtls" label="dtls"></el-option>
          <el-option value="wireguard" label="wireguard"></el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="路径" v-if="2 === parseInt(form.transport_type.toString())">
        <el-input v-model="form.web_socket.path" placeholder="URI路径"></el-input>
      </el-form-item>
      <el-form-item label="域名" v-if="4 === parseInt(form.transport_type.toString())">
        <el-input v-model="form.http2.host" placeholder="HTTP2的域名，多个使用英文,分隔"></el-input>
      </el-form-item>
      <el-form-item label="路径" v-if="4 === parseInt(form.transport_type.toString())">
        <el-input v-model="form.http2.path" placeholder="URI路径"></el-input>
      </el-form-item>
      <el-form-item label-width="0" class="content-center">
        <el-button type="danger" @click="close">取消</el-button>
        <el-button type="primary" @click="save" :loading="saving">保存</el-button>
      </el-form-item>
    </el-form>
  </el-dialog>
</template>

<script lang="ts">
import {defineComponent} from "vue"
import axios, {AxiosResponse} from "axios"
import {BaseResponse} from "@/api/base"
import {API_V2RAY_ENDPOINT_NEW, V2rayEndpointNewForm} from "@/api/v2ray_endpoint_new"

export default defineComponent({
  name: "NewV2rayEndpoint",

  props: {
    show: Boolean,
  },

  emits: ['update:show', 'success'],

  data() {
    return {
      saving: false,
      form: new V2rayEndpointNewForm(),
      rules: {
        cloud: [{required: true, message: '必须选择服务提供商', trigger: 'blur'}],
        endpoint: [{required: true, message: '必须选择节点所在位置', trigger: 'blur'}],
        host: [{required: true, message: '必须填写主机地址', trigger: 'blur'}],
        port: [{required: true, message: '必须填写端口号', trigger: 'blur'}],
        user_id: [{required: true, message: '必须填写用户ID', trigger: 'blur'}],
      },
    }
  },

  methods: {
    close() {
      this.form = new V2rayEndpointNewForm()

      this.$emit('update:show', false)
    },

    success() {
      this.close()

      this.$emit('success')
    },

    save() {
      let form = this.$refs.NewV2rayEndpointForm as any

      form.validate((isValid: boolean) => {
        if (!isValid) {
          return
        }

        this.form.port = parseInt(this.form.port.toString())
        this.form.alter_id = parseInt(this.form.alter_id.toString())

        this.saving = true
        axios.post(API_V2RAY_ENDPOINT_NEW, this.form).then((response: AxiosResponse<BaseResponse>) => {
          this.saving = false
          if (0 != response.data.code) {
            this.$message.error(response.data.message)

            return
          }

          this.$message.success('添加节点成功')
          this.success()
        })
      })
    },
  },
})
</script>

<style scoped>

</style>
