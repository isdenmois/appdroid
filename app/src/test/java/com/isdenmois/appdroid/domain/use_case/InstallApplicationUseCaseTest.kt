package com.isdenmois.appdroid.domain.use_case

import com.isdenmois.appdroid.domain.fake.ApkRepositoryFake
import com.isdenmois.appdroid.domain.model.AppPackage
import io.mockk.coEvery
import io.mockk.coVerify
import io.mockk.mockk
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.test.runTest
import org.junit.Test
import strikt.api.expectThat
import strikt.assertions.isEqualTo
import java.io.File

@ExperimentalCoroutinesApi
class InstallApplicationUseCaseTest {
    private val apkRepository = ApkRepositoryFake()
    private val downloadApk = mockk<DownloadApkUseCase>(relaxed = true)
    private val installApplication = InstallApplicationUseCase(apkRepository, downloadApk)

    @Test
    fun `should install app`() = runTest {
        val file = File("/tmp/app")
        val app = AppPackage(id = "test", appId = "com.test", type = "HZ", version = "1.0", versionName = "1.0")

        coEvery { downloadApk.invoke(any()) } returns file

        installApplication(app)

        coVerify { downloadApk.invoke(app) }
        expectThat(apkRepository.lastFile).isEqualTo(file)
    }
}
