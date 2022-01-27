package com.leo.architecturetraining.di

import com.leo.architecturetraining.data.HttpService
import com.leo.architecturetraining.data.datasource.CharacterDataSource
import com.leo.architecturetraining.data.datasource.CharacterDataSourceImpl
import com.leo.architecturetraining.data.repository.CharacterRepositoryImpl
import com.leo.architecturetraining.data.service.RickAndMortyApi
import com.leo.architecturetraining.domain.repository.CharacterRepository
import org.koin.androidx.viewmodel.dsl.viewModel
import org.koin.dsl.module

class CharacterModule {

    val characterModule = module {

        single<CharacterDataSource> { CharacterDataSourceImpl(
            service = HttpService().create(RickAndMortyApi::class)
            )
        }

        factory<CharacterRepository> { CharacterRepositoryImpl(
            characterDataSource = get()
            )
        }
    }
}