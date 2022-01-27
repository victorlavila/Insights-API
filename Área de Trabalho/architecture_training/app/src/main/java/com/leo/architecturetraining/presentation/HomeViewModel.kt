package com.leo.architecturetraining.presentation

import androidx.lifecycle.ViewModel
import com.leo.architecturetraining.domain.usecase.CharacterUseCase

class HomeViewModel(
    private val characterUseCase: CharacterUseCase
) : ViewModel() {
}