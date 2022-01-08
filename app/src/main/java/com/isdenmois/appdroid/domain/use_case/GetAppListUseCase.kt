package com.isdenmois.appdroid.domain.use_case

import com.isdenmois.appdroid.domain.model.AppPackage
import com.isdenmois.appdroid.domain.repository.AppPackagesRepository
import com.isdenmois.appdroid.domain.model.Resource
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.flow
import javax.inject.Inject

class GetAppListUseCase @Inject constructor(
    private val appPackagesRepository: AppPackagesRepository
) {
    operator fun invoke(): Flow<Resource<List<AppPackage>>> = flow {
        emit(Resource.loading(null))

        try {
            val appList = appPackagesRepository.fetchApps()

            emit(Resource.success(appList))
        } catch (exception: Exception) {
            exception.printStackTrace()
            emit(Resource.error(null, message = exception.message ?: "Error Occurred!"))
        }
    }
}
