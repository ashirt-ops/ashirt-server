import 'react-day-picker/dist/style.css';
import React, { useEffect, useState } from 'react';

import { DayPicker, SelectSingleEventHandler } from 'react-day-picker';
import classnames from 'classnames/bind';

import Popover from 'src/components/popover';
import Button from 'src/components/button';
import { setHours, setMinutes } from 'date-fns';

const cx = classnames.bind(require('./stylesheet'));

interface DateTimePickerProps {
	onSelectedDate: (date: Date) => void;
	selected?: Date | undefined;
}

interface PickerProps extends DateTimePickerProps {
	onButtonClick: () => void;
}

const Picker: React.FC<PickerProps> = ({ onButtonClick, onSelectedDate, selected }) => {
	const [date, setDate] = useState<Date>(selected ?? new Date());

	useEffect(() => {
		onSelectedDate(date)
	}, [date, onSelectedDate])

	const onChangeDate: SelectSingleEventHandler = (newDate) => {
		setDate(new Date(newDate as Date))
	}

	const onChangeHour: React.ChangeEventHandler<HTMLInputElement> = (event) => {
		const [hours, minutes] = event.target.value.split(":").map(Number)

		setDate(setMinutes(setHours(date, hours), minutes))
	}

	return (
		<div className={cx('popup')}>
			<div className={cx('day-picker-area')}>
				<DayPicker
					className={cx('day-picker')}
					mode="single"
					onSelect={onChangeDate}
					selected={date}
				/>
				<div className={cx("time-picker-wrapper")}>
					<input aria-label="Time" type="time" className={cx("time-picker")} onChange={onChangeHour} />
				</div>
				<Button primary className={cx('close-button')} onClick={onButtonClick}>
					Close
				</Button>
			</div>
		</div>
	);
};

const DateTimePicker: React.FC<DateTimePickerProps> = ({ onSelectedDate, selected }) => {
	const [isOpen, setIsOpen] = React.useState(false);

	return (
		<Popover
			isOpen={isOpen}
			onClick={() => setIsOpen(true)}
			onRequestClose={() => setIsOpen(false)}
			content={
				<Picker
					onButtonClick={() => setIsOpen(false)}
					onSelectedDate={onSelectedDate}
					selected={selected}
				/>
			}
		>
			<Button
				doNotSubmit
				className={cx('open-button')}
				icon={require('../date_range_picker/icon.svg')}
			/>
		</Popover>
	);
};

export default DateTimePicker;
