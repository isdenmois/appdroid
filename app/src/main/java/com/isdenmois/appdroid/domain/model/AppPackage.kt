package com.isdenmois.appdroid.domain.model

import androidx.compose.runtime.mutableStateOf
import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass

@JsonClass(generateAdapter = true)
data class AppPackage(
    val appId: String,
    var name: String? = null,
    val version: String,
    val versionName: String,
    var localVersion: String? = null,
    var localVersionName: String? = null,
) {
    @Json(ignore = true)
    val progress = mutableStateOf(-1)
}
