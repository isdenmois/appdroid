package com.isdenmois.appdroid.domain.fake

import com.isdenmois.appdroid.domain.model.AppPackage
import com.isdenmois.appdroid.domain.model.Download
import com.isdenmois.appdroid.domain.repository.AppPackagesRepository
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.asFlow

class AppPackagesRepositoryFake() : AppPackagesRepository {
    private var appsList: List<AppPackage> = listOf()
    private var downloads: List<Download> = listOf()
    private var exception: Exception? = null

    override suspend fun fetchApps(): List<AppPackage> {
        exception?.let {
            throw it
        }

        return appsList
    }

    override suspend fun downloadApp(app: AppPackage): Flow<Download> = downloads.asFlow()

    fun setList(vararg list: AppPackage) {
        appsList = list.asList()
        exception = null
    }

    fun setException(e: Exception) {
        exception = e
    }

    fun setDownloadList(vararg list: Download) {
        downloads = list.asList()
    }
}
