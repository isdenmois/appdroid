package com.isdenmois.appdroid.data.model

import com.google.gson.annotations.SerializedName

data class App(
    @SerializedName("id") val id: Int,
    @SerializedName("name") var name: String?,
    @SerializedName("appId") val appId: String,
    @SerializedName("version") val version: String,
    @SerializedName("type") val type: String,
    var localVersion: String?,
)
