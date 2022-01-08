package com.isdenmois.appdroid.data.di

import com.isdenmois.appdroid.data.repository.ApkRepositoryImpl
import com.isdenmois.appdroid.data.repository.AppPackagesRepositoryImpl
import com.isdenmois.appdroid.domain.repository.ApkRepository
import com.isdenmois.appdroid.domain.repository.AppPackagesRepository
import dagger.Binds
import dagger.Module
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent

@InstallIn(SingletonComponent::class)
@Module
abstract class DomainModule {
    @Binds
     abstract fun bindApkRepository(uploadRepository: ApkRepositoryImpl): ApkRepository

    @Binds
     abstract fun bindAppPackagesRepository(uploadRepository: AppPackagesRepositoryImpl): AppPackagesRepository
}
