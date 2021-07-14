package com.isdenmois.appdroid.shared.api

import retrofit2.http.GET

interface ApiService {
    @GET("/list")
    suspend fun getAppList(): List<AppPackage>
}
