package com.isdenmois.appdroid.domain.repository

import com.isdenmois.appdroid.domain.model.AppPackage
import com.isdenmois.appdroid.domain.model.Download
import kotlinx.coroutines.flow.Flow

interface AppPackagesRepository {
    suspend fun fetchApps(): List<AppPackage>

    suspend fun downloadApp(app: AppPackage): Flow<Download>
}
