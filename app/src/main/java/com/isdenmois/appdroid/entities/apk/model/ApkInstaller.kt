package com.isdenmois.appdroid.entities.apk.model

import android.app.DownloadManager
import android.content.*
import android.net.Uri
import android.os.Environment
import android.util.Log
import androidx.core.content.FileProvider
import com.isdenmois.appdroid.BuildConfig
import dagger.hilt.android.qualifiers.ApplicationContext
import kotlinx.coroutines.suspendCancellableCoroutine
import java.io.File
import javax.inject.Inject
import javax.inject.Named

class ApkInstaller @Inject constructor(
    @ApplicationContext private val context: Context,
    @Named("BASE_URL") private val baseUrl: String
) {
    suspend fun downloadAndInstallApplication(appId: String) {
        val apk = downloadApk(appId)

        installApplication(apk)
    }

    private suspend fun downloadApk(appId: String): File {
        val fileName = "$appId.apk"
        val url = "$baseUrl/$fileName"
        val file = File(context.getExternalFilesDir(Environment.DIRECTORY_DOWNLOADS), fileName)

        return suspendCancellableCoroutine { continuation ->
            if (file.exists()) {
                file.delete()
            }

            val request = DownloadManager.Request(Uri.parse(url))
                .setDestinationUri(Uri.fromFile(file)) // Uri of the destination file
                .setTitle(fileName) // Title of the Download Notification
                .setDescription("Downloading") // Description of the Download Notification
                .setRequiresCharging(false) // Set if charging is required to begin the download
                .setAllowedOverMetered(true) // Set if download is allowed on Mobile network
                .setAllowedOverRoaming(true) // Set if download is allowed on roaming network


            val downloadManager = context.getSystemService(Context.DOWNLOAD_SERVICE) as DownloadManager
            val downloadID = downloadManager.enqueue(request)
            val receiver = object : BroadcastReceiver() {
                override fun onReceive(context: Context, intent: Intent) {
                    val reference = intent.getLongExtra(DownloadManager.EXTRA_DOWNLOAD_ID, -1)

                    if (downloadID == reference) {
                        context.unregisterReceiver(this)
                        continuation.resume(file) {}
                    }
                }
            }

            continuation.invokeOnCancellation {
                context.unregisterReceiver(receiver)
            }

            context.registerReceiver(receiver, IntentFilter(DownloadManager.ACTION_DOWNLOAD_COMPLETE))
        }
    }

    private fun installApplication(file: File) {
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
