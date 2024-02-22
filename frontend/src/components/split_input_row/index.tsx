import classNames from 'classnames/bind';
import React from 'react';

import Input from 'src/components/input';
import WithLabel from 'src/components/with_label';

const cx = classNames.bind(require('./stylesheet'));

const SplitInputRow = (props: {
	label: string;
	inputValue: string;
	className?: string;
	children: React.ReactNode;
}) => (
	<WithLabel label={props.label}>
		<div className={cx('multi-item-row')}>
			<Input
				readOnly
				className={cx('flex-input', props.className)}
				value={props.inputValue}
			/>
			{props.children}
		</div>
	</WithLabel>
);

export default SplitInputRow;
