package com.isdenmois.appdroid.domain.use_case

import com.isdenmois.appdroid.domain.model.AppPackage
import com.isdenmois.appdroid.domain.model.Download
import com.isdenmois.appdroid.domain.repository.AppPackagesRepository
import kotlinx.coroutines.flow.lastOrNull
import kotlinx.coroutines.flow.mapNotNull
import java.io.File
import javax.inject.Inject

class DownloadApkUseCase @Inject constructor(
    private val appPackagesRepository: AppPackagesRepository,
    private val setProgress: SetProgressUseCase,
) {
    suspend operator fun invoke(app: AppPackage): File? {
        val result = appPackagesRepository.downloadApp(app).mapNotNull {
            setProgress(app, it)

            return@mapNotNull when(it) {
                is Download.Finished -> it.file
                else -> null
            }
        }.lastOrNull()

        return result
    }
}
