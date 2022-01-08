package com.isdenmois.appdroid.domain.use_case

import com.isdenmois.appdroid.domain.model.AppPackage
import com.isdenmois.appdroid.domain.repository.ApkRepository
import javax.inject.Inject

class InstallApplicationUseCase @Inject constructor(
    private val apkRepository: ApkRepository,
    private val downloadApk: DownloadApkUseCase,
) {
    suspend operator fun invoke(app: AppPackage) {
        val apk = downloadApk(app)

        apk?.let {
            apkRepository.installApplication(apk)
        }
    }
}
