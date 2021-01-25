package com.isdenmois.appdroid.utils

import android.content.ActivityNotFoundException
import android.content.Context
import android.content.Intent
import android.net.Uri
import android.util.Log
import androidx.core.content.FileProvider
import com.isdenmois.appdroid.BuildConfig
import java.io.File


object ApkInstaller {
    fun installApplication(context: Context, file: File) {
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
        return FileProvider.getUriForFile(context, BuildConfig.APPLICATION_ID.toString() + ".provider", file)
    }
}