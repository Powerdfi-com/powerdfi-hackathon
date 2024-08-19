"use client"
import React, { ReactNode } from 'react'
import Slider from 'react-slick';

const CustomSlider = ({ children, className, slidesToShow }: { children: ReactNode, className?: string; slidesToShow?: number }) => {
    const settings = {
        dots: true,
        infinite: true,
        speed: 500,
        slidesToShow: slidesToShow ?? 1,
        slidesToScroll: 1,

    };
    return <Slider {...settings} className={className}>
        {children}
    </Slider >;
}

export default CustomSlider