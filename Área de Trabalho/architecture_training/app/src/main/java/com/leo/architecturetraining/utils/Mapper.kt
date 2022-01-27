package com.leo.architecturetraining.utils

interface Mapper<S, T> {
    fun map(source: S): T
}