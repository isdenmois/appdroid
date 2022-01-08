package com.isdenmois.appdroid.domain.use_case

import com.isdenmois.appdroid.domain.fake.AppPackagesRepositoryFake
import com.isdenmois.appdroid.domain.model.AppPackage
import com.isdenmois.appdroid.domain.model.Status
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.flow.toList
import kotlinx.coroutines.test.runTest
import org.junit.Test
import strikt.api.expectThat
import strikt.assertions.isEqualTo
import strikt.assertions.isNotNull
import strikt.assertions.isNull

@ExperimentalCoroutinesApi
class GetAppListUseCaseTest {
    private val appPackagesRepository = AppPackagesRepositoryFake()
    private val getAppList = GetAppListUseCase(appPackagesRepository)

    @Test
    fun `should return app list`() = runTest {
        appPackagesRepository.setList(
            AppPackage(id = "test", appId = "com.test", type = "HZ", version = "1.0", versionName = "1.0"),
            AppPackage(id = "haro", appId = "com.haro", type = "HZ", version = "1.0", versionName = "1.0"),
            AppPackage(id = "hello", appId = "com.hello", type = "HZ", version = "1.0", versionName = "1.0"),
        )

        val emits = getAppList().toList()

        expectThat(emits.size).isEqualTo(2)
        expectThat(emits[0].status).isEqualTo(Status.LOADING)

        expectThat(emits[1].status).isEqualTo(Status.SUCCESS)
        expectThat(emits[1].data).isNotNull()
        expectThat(emits[1].data!!.size).isEqualTo(3)
    }

    @Test
    fun `should return error`() = runTest {
        appPackagesRepository.setException(Exception("Something went right"))

        val emits = getAppList().toList()

        expectThat(emits.size).isEqualTo(2)
        expectThat(emits[0].status).isEqualTo(Status.LOADING)

        expectThat(emits[1].status).isEqualTo(Status.ERROR)
        expectThat(emits[1].data).isNull()
        expectThat(emits[1].message).isEqualTo("Something went right")
    }
}
