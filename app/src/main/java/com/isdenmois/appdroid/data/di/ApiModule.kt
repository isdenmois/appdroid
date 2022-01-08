package com.isdenmois.appdroid.data.di

import android.content.Context
import com.isdenmois.appdroid.R
import com.isdenmois.appdroid.data.remote.ApiService
import com.squareup.moshi.Moshi
import com.squareup.moshi.kotlin.reflect.KotlinJsonAdapterFactory
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.android.qualifiers.ApplicationContext
import dagger.hilt.components.SingletonComponent
import okhttp3.OkHttpClient
import retrofit2.Retrofit
import retrofit2.converter.moshi.MoshiConverterFactory
import retrofit2.converter.scalars.ScalarsConverterFactory
import javax.inject.Named
import javax.inject.Singleton

@InstallIn(SingletonComponent::class)
@Module
object ApiModule {
    @Provides
    @Named("BASE_URL")
    fun provideBaseURL(@ApplicationContext context: Context) = context.getString(R.string.base_url)

    @Singleton
    @Provides
    fun provideRetrofitService(@Named("BASE_URL") baseUrl: String): ApiService {
        val builder = OkHttpClient.Builder()

        val okHttpClient = builder.build()
        val moshi = Moshi.Builder()
            .add(KotlinJsonAdapterFactory())
            .build()

        return Retrofit.Builder()
            .baseUrl(baseUrl)
            .addConverterFactory(ScalarsConverterFactory.create())
            .addConverterFactory(MoshiConverterFactory.create(moshi))
            .client(okHttpClient)
            .build()
            .create(ApiService::class.java)
    }
}
