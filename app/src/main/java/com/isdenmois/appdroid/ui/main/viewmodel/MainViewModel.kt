package com.isdenmois.appdroid.ui.main.viewmodel

import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import com.isdenmois.appdroid.data.model.App
import com.isdenmois.appdroid.data.repository.MainRepository
import com.isdenmois.appdroid.utils.Resource
import io.reactivex.android.schedulers.AndroidSchedulers
import io.reactivex.disposables.CompositeDisposable
import io.reactivex.schedulers.Schedulers

class MainViewModel(private val mainRepository: MainRepository) : ViewModel() {
    private val apps = MutableLiveData<Resource<Array<App>>>()
    private val compositeDisposable = CompositeDisposable()

    init {
        fetchApps()
    }

    override fun onCleared() {
        super.onCleared()
        compositeDisposable.dispose()
    }

    fun getApps(): LiveData<Resource<Array<App>>> = apps

    fun fetchApps() {
        apps.postValue(Resource.loading(null))

        compositeDisposable.add(
            mainRepository.getApps()
                .map { list -> list.map { app ->
                    val packageInfo = mainRepository.getPackage(app.appId)
                    if (packageInfo != null) {
                        app.name = packageInfo.packageName
                        app.localVersion = packageInfo.versionCode.toString()
                    }

                    app
                }}
                .subscribeOn(Schedulers.io())
                .observeOn(AndroidSchedulers.mainThread())
                .subscribe({ appList ->
                    apps.postValue(Resource.success(appList.toTypedArray()))
                }, { throwable ->
                    apps.postValue(Resource.error("Something Went Wrong", null))
                })
        )
    }
}
