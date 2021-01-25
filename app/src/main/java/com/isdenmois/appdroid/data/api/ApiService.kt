package com.isdenmois.appdroid.data.api

import android.app.DownloadManager
import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent
import android.content.IntentFilter
import android.content.pm.PackageInfo
import android.net.Uri
import android.os.Environment
import androidx.appcompat.app.AppCompatActivity
import com.isdenmois.appdroid.data.model.App
import com.rx2androidnetworking.Rx2AndroidNetworking
import io.reactivex.Single
import java.io.File

interface OnApkDownloadedListener {
    fun onSuccess(file: File)
}

class ApiService(val context: Context) {
    private val baseUrl = "http://192.168.1.60:3000"
//    private val baseUrl = "http://192.168.1.60:3000"

    fun getApps(): Single<List<App>> {
        return Rx2AndroidNetworking.get("$baseUrl/list")
            .build()
            .getObjectListSingle(App::class.java)
    }

    fun getPackage(appId: String): PackageInfo? {
        val packages = context.packageManager.getInstalledPackages(0)

        return packages.firstOrNull { p -> p.packageName.equals(appId)}
    }

    fun getApk(appId: String, listener: OnApkDownloadedListener) {
        val fileName = "$appId.apk"
        val url = "$baseUrl/$fileName"
        val file = File(context.getExternalFilesDir(Environment.DIRECTORY_DOWNLOADS), fileName)

        val request = DownloadManager.Request(Uri.parse(url))
            .setDestinationUri(Uri.fromFile(file)) // Uri of the destination file
            .setTitle(fileName) // Title of the Download Notification
            .setDescription("Downloading") // Description of the Download Notification
            .setRequiresCharging(false) // Set if charging is required to begin the download
            .setAllowedOverMetered(true) // Set if download is allowed on Mobile network
            .setAllowedOverRoaming(true) // Set if download is allowed on roaming network


        val downloadManager = context.getSystemService(AppCompatActivity.DOWNLOAD_SERVICE) as DownloadManager
        val downloadID = downloadManager.enqueue(request)
        val receiver = object : BroadcastReceiver() {
            override fun onReceive(context: Context, intent: Intent) {
                val reference = intent.getLongExtra(DownloadManager.EXTRA_DOWNLOAD_ID, -1)

                if (downloadID == reference) {
                    context.unregisterReceiver(this)
                    listener.onSuccess(file)
                }
            }
        }

        context.registerReceiver(receiver, IntentFilter(DownloadManager.ACTION_DOWNLOAD_COMPLETE))
    }
}
