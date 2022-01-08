package com.isdenmois.appdroid.shared.api

import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.distinctUntilChanged
import kotlinx.coroutines.flow.flow
import kotlinx.coroutines.flow.flowOn
import okhttp3.ResponseBody
import java.io.File

sealed class Download {
    data class Progress(val percent: Int) : Download()
    data class Finished(val file: File) : Download()
}

fun ResponseBody.downloadToFileWithProgress(directory: File, filename: String): Flow<Download> = flow {
    emit(Download.Progress(0))

    // flag to delete file if download errors or is cancelled
    var deleteFile = true
    val file = File(directory, "${filename}.${contentType()?.subtype}")

    try {
        byteStream().use { inputStream ->
            file.outputStream().use { outputStream ->
                val totalBytes = contentLength()
                val data = ByteArray(8_192)
                var progressBytes = 0L

                while (true) {
                    val bytes = inputStream.read(data)

                    if (bytes == -1) {
                        break
                    }

                    outputStream.channel
                    outputStream.write(data, 0, bytes)
                    progressBytes += bytes

                    emit(Download.Progress(percent = ((progressBytes * 100) / totalBytes).toInt()))
                }

                when {
                    progressBytes < totalBytes ->
                        throw Exception("missing bytes")
                    progressBytes > totalBytes ->
                        throw Exception("too many bytes")
                    else ->
                        deleteFile = false
                }
            }
        }

        emit(Download.Finished(file))
    } finally {
        // check if download was successful

        if (deleteFile) {
            file.delete()
        }
    }
}
    .flowOn(Dispatchers.IO)
    .distinctUntilChanged()
