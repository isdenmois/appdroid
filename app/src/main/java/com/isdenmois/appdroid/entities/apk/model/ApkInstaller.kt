package com.isdenmois.appdroid.entities.apk.model

import android.content.ActivityNotFoundException
import android.content.Context
import android.content.Intent
import android.net.Uri
import android.util.Log
import androidx.core.content.FileProvider
import com.isdenmois.appdroid.BuildConfig
import dagger.hilt.android.qualifiers.ApplicationContext
import java.io.File
import javax.inject.Inject

class ApkInstaller @Inject constructor(
    @ApplicationContext private val context: Context,
) {
    fun installApplication(file: File) {
        val intent = Intent(Intent.ACTION_VIEW)

        intent.setDataAndType(uriFromFile(context, file), "application/vnd.android.package-archive")
        intent.addFlags(Intent.FLAG_ACTIVITY_NEW_TASK)
        intent.addFlags(Intent.FLAG_GRANT_READ_URI_PERMISSION)

        try {
            context.startActivity(intent)
        } catch (e: ActivityNotFoundException) {
            e.printStackTrace()
            Log.e("TAG", "Error in opening the file!")
        }
    }

    private fun uriFromFile(context: Context, file: File): Uri {
        return FileProvider.getUriForFile(context, BuildConfig.APPLICATION_ID + ".provider", file)
    }
}
