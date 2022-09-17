package com.isdenmois.appdroid.data.remote

import com.isdenmois.appdroid.domain.model.AppPackage
import okhttp3.ResponseBody
import retrofit2.http.GET
import retrofit2.http.Path
import retrofit2.http.Streaming

interface ApiService {
    @GET("list")
    suspend fun getAppList(): List<AppPackage>

    @Streaming
    @GET("/{fileName}")
    suspend fun getApk(@Path("fileName") fileName: String): ResponseBody
}
